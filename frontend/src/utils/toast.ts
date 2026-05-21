let _toast: any = null

export function setToast(fn: any) { _toast = fn }

export function toast(msg: string, type = 'info') {
  if (_toast) _toast(msg, type)
  else alert(msg)
}
