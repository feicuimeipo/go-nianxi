interface Entity {
    title: string;
    log(): void;
  }
 
  class Post implements Entity {
    title: string;
 
    constructor(title: string) {
      this.title = title;
    }
 
    log(): void {
      console.log(this.title);
    }
  }