package ansible

type Schedule struct {
	ID             int     `db:"id" json:"id"`
	ProjectID      int     `db:"project_id" json:"project_id"`
	TemplateID     int     `db:"template_id" json:"template_id"`
	CronFormat     string  `db:"cron_format" json:"cron_format"`
	RepositoryID   *int    `db:"repository_id" json:"repository_id"`
	LastCommitHash *string `db:"last_commit_hash" json:"-"`
}
