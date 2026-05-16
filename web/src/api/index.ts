import api from '@/utils/http'
import type {
	LoginRequest,
	RoleRequest, RoleResponse,
	PaginationRequest, PaginationResponse,
	CreateMenuRequest, UpdateMenuRequest, MenuResponse, MenuTreeNode,
	CreateOrgUnitRequest, UpdateOrgUnitRequest, OrgUnitResponse,
	ResourceResponse,
	BindResourceRequest, BindMenuRequest, BindMenuResourceRequest,
	UserPermissions,
	UserResponse, CreateUserRequest, UpdateUserRequest,
	ProfileResponse, UpdateProfileRequest, ChangePasswordRequest,
	DashboardStats, AuditTrendItem,
	RoleDetailResponse,
	DictResponse, DictDetailResponse,
	CreateDictRequest, UpdateDictRequest,
	CreateDictDetailRequest, UpdateDictDetailRequest,
} from '@/types'

export const login = (data: LoginRequest) => api.post('/login', data)

export const getRoles = (params?: PaginationRequest) => api.get<{ data: PaginationResponse<RoleResponse> }>('/roles', { params })
export const getRoleDetail = (id: string) => api.get<{ data: RoleDetailResponse }>(`/roles/${id}`)
export const createRole = (data: RoleRequest) => api.post<{ data: RoleResponse }>('/roles', data)
export const updateRole = (id: string, data: RoleRequest) => api.put(`/roles/${id}`, data)
export const deleteRole = (id: string) => api.delete(`/roles/${id}`)
export const bindRoleResource = (roleId: string, resourceId: string, data: BindResourceRequest) => api.post(`/roles/${roleId}/resources`, { ...data, resource_id: resourceId })
export const unbindRoleResource = (roleId: string, resourceId: string) => api.delete(`/roles/${roleId}/resources/${resourceId}`)
export const bindRoleMenu = (roleId: string, data: BindMenuRequest) => api.post(`/roles/${roleId}/menus`, data)
export const unbindRoleMenu = (roleId: string, menuId: string) => api.delete(`/roles/${roleId}/menus/${menuId}`)

export const getMenuTree = () => api.get<{ data: MenuTreeNode[] }>('/menus/tree')
export const getMenu = (id: string) => api.get<{ data: MenuResponse }>(`/menus/${id}`)
export const createMenu = (data: CreateMenuRequest) => api.post('/menus', data)
export const updateMenu = (id: string, data: UpdateMenuRequest) => api.put(`/menus/${id}`, data)
export const deleteMenu = (id: string) => api.delete(`/menus/${id}`)
export const bindMenuResource = (menuId: string, data: BindMenuResourceRequest) => api.post(`/menus/${menuId}/resources`, data)
export const unbindMenuResource = (menuId: string, resourceId: string) => api.delete(`/menus/${menuId}/resources/${resourceId}`)

export const getOrgTree = () => api.get<{ data: OrgUnitResponse[] }>('/org-units/tree')
export const createOrgUnit = (data: CreateOrgUnitRequest) => api.post('/org-units', data)
export const updateOrgUnit = (id: string, data: UpdateOrgUnitRequest) => api.put(`/org-units/${id}`, data)
export const deleteOrgUnit = (id: string) => api.delete(`/org-units/${id}`)

export const getResources = (params?: PaginationRequest) => api.get<{ data: PaginationResponse<ResourceResponse> }>('/resources', { params })
export const createResource = (data: ResourceResponse) => api.post('/resources', data)
export const updateResource = (id: string, data: ResourceResponse) => api.put(`/resources/${id}`, data)
export const deleteResource = (id: string) => api.delete(`/resources/${id}`)

export const getAuditLogs = (params?: PaginationRequest) => api.get<{ data: PaginationResponse<any> }>('/audit-logs', { params })
export const getAuditLogsByTarget = (params: { target_type: string; target_id: string }) => api.get<{ data: { data: any[] } }>('/audit-logs/target', { params })
export const getAuditLogsByTime = (params: { start_time: string; end_time: string }) => api.get<{ data: { data: any[] } }>('/audit-logs/time', { params })

export const getUserPermissions = () => api.get<{ data: UserPermissions }>('/user/permissions')
export const getUserMenus = () => api.get<{ data: { menus: MenuTreeNode[] } }>('/user/menus')
export const getProfile = () => api.get<{ data: ProfileResponse }>('/user/profile')
export const updateProfile = (data: UpdateProfileRequest) => api.put('/user/profile', data)
export const changePassword = (data: ChangePasswordRequest) => api.put('/user/password', data)

export const getUsers = (params?: PaginationRequest) => api.get<{ data: PaginationResponse<UserResponse> }>('/users', { params })
export const getUser = (id: string) => api.get<{ data: UserResponse }>(`/users/${id}`)
export const createUser = (data: CreateUserRequest) => api.post('/users', data)
export const updateUser = (id: string, data: UpdateUserRequest) => api.put(`/users/${id}`, data)
export const deleteUser = (id: string) => api.delete(`/users/${id}`)

export const getDashboardStats = () => api.get<{ data: DashboardStats }>('/dashboard/stats')
export const getAuditTrend = () => api.get<{ data: AuditTrendItem[] }>('/dashboard/audit-trend')

export const getDicts = (params?: PaginationRequest) => api.get<{ data: PaginationResponse<DictResponse> }>('/dict', { params })
export const getDict = (id: string) => api.get<{ data: DictResponse }>(`/dict/${id}`)
export const createDict = (data: CreateDictRequest) => api.post('/dict', data)
export const updateDict = (id: string, data: UpdateDictRequest) => api.put(`/dict/${id}`, data)
export const deleteDict = (id: string) => api.delete(`/dict/${id}`)
export const getDictDetails = (id: string) => api.get<{ data: DictDetailResponse[] }>(`/dict/${id}/details`)
export const createDictDetail = (dictId: string, data: CreateDictDetailRequest) => api.post(`/dict/${dictId}/details`, data)
export const updateDictDetail = (dictId: string, detailId: string, data: UpdateDictDetailRequest) => api.put(`/dict/${dictId}/details/${detailId}`, data)
export const deleteDictDetail = (dictId: string, detailId: string) => api.delete(`/dict/${dictId}/details/${detailId}`)

// 公共字典查询：根据类型获取字典值列表
export const getDictInfo = (type: string) => api.get<{ data: DictDetailResponse[] }>(`/dict-info/${type}`)
