export default class Cookies {
  /**
   * 获取cookie
   * @param key cookie名称/键名
   */
  static get(key: string): string {
    const name = key + "=";
    const ca = document.cookie.split(";");
    for (let i = 0; i < ca.length; i++) {
      const c = ca[i].trim();
      if (c.indexOf(name) === 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
  }

  /**
   * 设置cookie
   * @param key cookie名称/键名
   * @param value cookie值
   * @param exp cookie到期时间戳
   */
  static set(key: string, value: string) {
    let exp = this.getTime(-13)
    const expires = "expires=" + new Date(exp).toUTCString();
    return (document.cookie = key + "=" + value + "; " + expires);
  }

  /**
   * 删除cookie
   * @param key cookie名称/键名
   */
  static remove(key: string) {
    return (document.cookie = `${key}=; expires=Thu, 01 Jan 1970 00:00:00 GMT`);
  }


  /**
   *     //上周的开始时间
   *     console.log(getTime(7));
   *     //上周的结束时间
   *     console.log(getTime(1));
   *     //本周的开始时间
   *     console.log(getTime(0));
   *     //本周的结束时间
   *     console.log(getTime(-6));
   *     //下周的开始时间
   *     console.log(getTime(-7));
   *     //下周结束时间
   *     console.log(getTime(-13));
   * @param n
   */
  static getTime(n:number) {
    let now = new Date();
    let year = now.getFullYear();
    let month = now.getMonth() + 1;
    let day = now.getDay(); //返回星期几的某一天;
    n = day === 0 ? n + 6 : n + (day - 1)
    now.setDate(now.getDate() - n);
    let date = now.getDate();
    let s = year + "-" + (month < 10 ? ('0' + month) : month) + "-" + (date < 10 ? ('0' + date) : date);
    return s;
  }

}
