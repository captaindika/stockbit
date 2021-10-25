SELECT u.ID, u.UserName, parent.name AS ParentUserName FROM USER u 
	LEFT JOIN (SELECT ID, UserName FROM USER)parent ON parent.id = u.parent