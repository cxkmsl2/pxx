export {}
declare global {
  interface Window {
    $toast: (msg: string, type?: string) => void
  }
}
