/** Minimal ZIP writer (STORE/deflate via CRC32) for browser exports */

const CRC_TABLE = (() => {
  const t = new Uint32Array(256);
  for (let i = 0; i < 256; i++) {
    let c = i;
    for (let k = 0; k < 8; k++) c = (c & 1) ? (0xedb88320 ^ (c >>> 1)) : (c >>> 1);
    t[i] = c >>> 0;
  }
  return t;
})();

function crc32(buf) {
  let c = 0xffffffff;
  for (let i = 0; i < buf.length; i++) c = CRC_TABLE[(c ^ buf[i]) & 0xff] ^ (c >>> 8);
  return (c ^ 0xffffffff) >>> 0;
}

function u16(n) { return [n & 0xff, (n >>> 8) & 0xff]; }
function u32(n) { return [n & 0xff, (n >>> 8) & 0xff, (n >>> 16) & 0xff, (n >>> 24) & 0xff]; }

function encodeName(name) {
  return new TextEncoder().encode(name.replace(/\\/g, '/'));
}

/**
 * @param {Record<string, string|Uint8Array>} files path -> content
 * @returns {Blob}
 */
export function createZipBlob(files) {
  const parts = [];
  const central = [];
  let offset = 0;

  for (const [name, content] of Object.entries(files)) {
    const data = typeof content === 'string' ? new TextEncoder().encode(content) : content;
    const nameBytes = encodeName(name);
    const crc = crc32(data);
    const size = data.length;

    const local = new Uint8Array([
      0x50, 0x4b, 0x03, 0x04, 20, 0, 0, 0, 0, 0,
      ...u32(crc), ...u32(size), ...u32(size),
      ...u16(nameBytes.length), 0, 0,
      ...nameBytes, ...data
    ]);
    parts.push(local);
    central.push({ nameBytes, crc, size, offset });
    offset += local.length;
  }

  const centralStart = offset;
  for (const e of central) {
    const cd = new Uint8Array([
      0x50, 0x4b, 0x01, 0x02, 20, 0, 20, 0, 0, 0, 0, 0,
      ...u32(e.crc), ...u32(e.size), ...u32(e.size),
      ...u16(e.nameBytes.length), 0, 0, 0, 0, 0, 0,
      ...u32(e.offset), ...e.nameBytes
    ]);
    parts.push(cd);
    offset += cd.length;
  }

  const count = central.length;
  const end = new Uint8Array([
    0x50, 0x4b, 0x05, 0x06, 0, 0, 0, 0,
    ...u16(count), ...u16(count), ...u32(offset - centralStart), ...u32(centralStart), 0, 0
  ]);
  parts.push(end);

  return new Blob(parts, { type: 'application/zip' });
}

export function downloadBlob(blob, filename) {
  const a = document.createElement('a');
  a.href = URL.createObjectURL(blob);
  a.download = filename;
  a.click();
  setTimeout(() => URL.revokeObjectURL(a.href), 2000);
}
