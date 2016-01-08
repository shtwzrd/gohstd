export class NewPost {
    constructor(
        public title :string,
        public message :string
    ){}
}

export class Post {
    constructor(
        public title :string,
        public message :string,
        public username :string,
        public timestamp :string
	  ){}
}
