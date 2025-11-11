
// src/utils/cookies.js
export function setCookie(name, value, days=7, path='/', sameSite='Lax') {
  const d = new Date();
  d.setTime(d.getTime() + (days*24*60*60*1000));
  const expires = "expires=" + d.toUTCString();
  document.cookie = `${name}=${encodeURIComponent(value)};${expires};path=${path};SameSite=${sameSite}`;
}
export function getCookie(name) {
  const cname = name + "=";
  const decodedCookie = decodeURIComponent(document.cookie);
  const ca = decodedCookie.split(';');
  for(let c of ca){
    c = c.trim();
    if (c.indexOf(cname) === 0) return c.substring(cname.length, c.length);
  }
  return null;
}
export function eraseCookie(name, path='/') {
  document.cookie = name+'=; Max-Age=-99999999; path=' + path;
}
