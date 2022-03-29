package master

func (m *master) fillConfig() {
	t := m.masterConfig.TemporaryFields()
	t = append(t, "password")
	m.masterConfig.SetTemporaryFields(t)

	d := m.masterConfig.DefaultValues()
	d["file"] = "~/.pass-keeper.master"
	m.masterConfig.SetDefaultValues(d)

	i := m.masterConfig.InstallFields()
	i = append(i, "file")
	m.masterConfig.SetInstallFields(i)

	f := m.masterConfig.FieldNames()
	f["file"] = "файл с мастер паролем"
	m.masterConfig.SetFieldNames(f)

	m.masterConfig.SetInit(m.validateConfig)
}
