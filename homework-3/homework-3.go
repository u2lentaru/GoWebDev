package main

func main() {
	type TPost struct {
		Subj     string
		PostTime string
		Text     string
	}

	type TBlog struct {
		Title    string
		PostList []TPost
	}

	var MyBlog = TBlog{
		Title: "My blog",
		PostList: []TPost{
			TPost{"1st subj", "01.01.2020", "1st text"},
			TPost{"2nd subj", "02.01.2020", "2nd text"},
			TPost{"3rd subj", "03.01.2020", "3rd text"},
		},
	}

}
