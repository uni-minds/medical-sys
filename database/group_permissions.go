package database

const (
	PermissionUserMediaList     = 0x01
	PermissionUserMediaManage   = 0x02
	PermissionUserUsersList     = 0x04
	PermissionUserUsersManage   = 0x06
	PermissionUserLabelsList    = 0x08
	PermissionUserLabelsManage  = 0x10
	PermissionUserReviewsList   = 0x20
	PermissionUserReviewsManage = 0x40
	PermissionUserGroupMaster   = PermissionUserMediaList | PermissionUserMediaManage | PermissionUserUsersList | PermissionUserUsersManage | PermissionUserLabelsList | PermissionUserLabelsManage | PermissionUserReviewsList | PermissionUserReviewsManage
)

func setPermissionsRole(role string) (v int) {
	p := UserPermissions{}

	p.ListMedia = true

	if role == "member" || role == "leader" || role == "master" {
		p.ListLabels = true
		p.ListUsers = true
	}

	if role == "leader" || role == "master" {
		p.ManageReviews = true
		p.ListReviews = true
	}

	if role == "master" {
		p.ManageLabels = true
		p.ManageReviews = true
		p.ManageUsers = true
	}
	return setPermissions(p)
}

func setPermissions(permissions UserPermissions) (p int) {
	if permissions.ListMedia {
		p |= PermissionUserMediaList
	}

	if permissions.ManageMedia {
		p |= PermissionUserMediaManage
	}

	if permissions.ListUsers {
		p |= PermissionUserUsersList
	}

	if permissions.ManageUsers {
		p |= PermissionUserUsersManage
	}

	if permissions.ListLabels {
		p |= PermissionUserLabelsList
	}

	if permissions.ManageLabels {
		p |= PermissionUserLabelsManage
	}

	if permissions.ListReviews {
		p |= PermissionUserReviewsList
	}

	if permissions.ManageReviews {
		p |= PermissionUserReviewsManage
	}
	return
}
func getPermissions(p int) (permissions UserPermissions) {
	if p == 0 {
		return
	}

	return UserPermissions{
		ListMedia:     p&PermissionUserMediaList != 0,
		ManageMedia:   p&PermissionUserMediaManage != 0,
		ListUsers:     p&PermissionUserUsersList != 0,
		ManageUsers:   p&PermissionUserUsersManage != 0,
		ListLabels:    p&PermissionUserLabelsList != 0,
		ManageLabels:  p&PermissionUserLabelsManage != 0,
		ListReviews:   p&PermissionUserReviewsList != 0,
		ManageReviews: p&PermissionUserReviewsManage != 0,
	}
}
