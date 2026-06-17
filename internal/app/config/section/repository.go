package section

type (
	Repository struct {
		Postgres RepositoryPostgres
	}

	RepositoryPostgres struct {
		Address  string `required:"true" default:"localhost" split_words:"true"`
		Username string `required:"true" split_words:"true"`
		Password string `required:"true" split_words:"true"`
		Name     string `required:"true" split_words:"true"`
	}
)
