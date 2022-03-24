package ansible

import "github.com/flipped-aurora/gin-vue-admin/server/model/ansible"

type TaskService struct {
}

func (taskService *TaskService) GetIncomingVersion(task *ansible.Task) *string {
	if task.BuildTaskID == nil {
		return nil
	}

	buildTask, err := d.GetTask(task.ProjectID, *task.BuildTaskID)

	if err != nil {
		return nil
	}

	tpl, err := d.GetTemplate(task.ProjectID, buildTask.TemplateID)
	if err != nil {
		return nil
	}

	if tpl.Type == TemplateBuild {
		return buildTask.Version
	}

	return buildTask.GetIncomingVersion(d)
}

func (task *Task) ValidateNewTask(template Template) error {
	switch template.Type {
	case TemplateBuild:
	case TemplateDeploy:
	case TemplateTask:
	}
	return nil
}

func (task *TaskWithTpl) Fill(d Store) error {
	if task.BuildTaskID != nil {
		build, err := d.GetTask(task.ProjectID, *task.BuildTaskID)
		if err == ErrNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		task.BuildTask = &build
	}
	return nil
}

