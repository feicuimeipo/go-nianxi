import CryptoJS from 'crypto-js'

export const aesEncrypt = (word, keyWord="XwKsGlMcdPMEhR1B") => {
  const key = CryptoJS.enc.Utf8.parse(keyWord);
  const src = CryptoJS.enc.Utf8.parse(word);
  const encrypted = CryptoJS.AES.encrypt(src, key, {mode:CryptoJS.mode.ECB,padding: CryptoJS.pad.Pkcs7});
  return encrypted.toString();
}
