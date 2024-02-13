export function getNonceFromFileName(fileName: string): string {
  return fileName.substring(0, 12);
}
