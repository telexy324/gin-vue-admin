import { getProjectList } from '@/api/ansibleProject'

export default {
  data() {
    return {
      projects: [],
      currentProject: {}
    }
  },
  methods: {
    async getProjects() {
      const table = await getProjectList({ page: 1, pageSize: 10 })
      if (table.code === 0) {
        this.projects = table.data.list
        this.currentProject = this.projects[0]
        this.searchInfo = { projectId: this.currentProject.ID }
      }
    },
    async setCurrentProject(val) {
      this.currentProject = val
      this.searchInfo = { projectId: this.currentProject.ID }
      await this.getTableData()
    }
  }
}
