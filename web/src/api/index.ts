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
  RoleDetailResponse
} from '@/types'

export const login = (data: LoginRequest) => api.post('/login', data)

export const getRoles = (params?: PaginationRequest) => api.get<{ data: PaginationResponse<RoleResponse> }>('/roles', { params })
export const getRoleDetail = (id: number) => api.get<{ data: RoleDetailResponse }>(`/roles/${id}`)
export const createRole = (data: RoleRequest) => api.post<{ data: RoleResponse }>('/roles', data)
export const updateRole = (id: number, data: RoleRequest) => api.put(`/roles/${id}`, data)
export const deleteRole = (id: number) => api.delete(`/roles/${id}`)
export const bindRoleResource = (roleId: number, resourceId: number, data: BindResourceRequest) => api.post(`/roles/${roleId}/resources`, { ...data, resource_id: resourceId })
export const unbindRoleResource = (roleId: number, resourceId: number) => api.delete(`/roles/${roleId}/resources/${resourceId}`)
export const bindRoleMenu = (roleId: number, data: BindMenuRequest) => api.post(`/roles/${roleId}/menus`, data)
export const unbindRoleMenu = (roleId: number, menuId: number) => api.delete(`/roles/${roleId}/menus/${menuId}`)

export const getMenuTree = () => api.get<{ data: MenuTreeNode[] }>('/menus/tree')
export const getMenu = (id: number) => api.get<{ data: MenuResponse }>(`/menus/${id}`)
export const createMenu = (data: CreateMenuRequest) => api.post('/menus', data)
export const updateMenu = (id: number, data: UpdateMenuRequest) => api.put(`/menus/${id}`, data)
export const deleteMenu = (id: number) => api.delete(`/menus/${id}`)
export const bindMenuResource = (menuId: number, data: BindMenuResourceRequest) => api.post(`/menus/${menuId}/resources`, data)
export const unbindMenuResource = (menuId: number, resourceId: number) => api.delete(`/menus/${menuId}/resources/${resourceId}`)

export const getOrgTree = () => api.get<{ data: OrgUnitResponse[] }>('/org-units/tree')
export const createOrgUnit = (data: CreateOrgUnitRequest) => api.post('/org-units', data)
export const updateOrgUnit = (id: number, data: UpdateOrgUnitRequest) => api.put(`/org-units/${id}`, data)
export const deleteOrgUnit = (id: number) => api.delete(`/org-units/${id}`)

export const getResources = (params?: PaginationRequest) => api.get<{ data: PaginationResponse<ResourceResponse> }>('/resources', { params })
export const createResource = (data: ResourceResponse) => api.post('/resources', data)
export const updateResource = (id: number, data: ResourceResponse) => api.put(`/resources/${id}`, data)
export const deleteResource = (id: number) => api.delete(`/resources/${id}`)

export const getAuditLogs = (params?: PaginationRequest) => api.get<{ data: PaginationResponse<any> }>('/audit-logs', { params })
export const getAuditLogsByTarget = (params: { target_type: string; target_id: string }) => api.get<{ data: { data: any[] } }>('/audit-logs/target', { params })
export const getAuditLogsByTime = (params: { start_time: string; end_time: string }) => api.get<{ data: { data: any[] } }>('/audit-logs/time', { params })

export const getUserPermissions = () => api.get<{ data: UserPermissions }>('/user/permissions')
export const getUserMenus = () => api.get<{ data: { menus: MenuTreeNode[] } }>('/user/menus')
export const getProfile = () => api.get<{ data: ProfileResponse }>('/user/profile')
export const updateProfile = (data: UpdateProfileRequest) => api.put('/user/profile', data)
export const changePassword = (data: ChangePasswordRequest) => api.put('/user/password', data)

export const getUsers = (params?: PaginationRequest) => api.get<{ data: PaginationResponse<UserResponse> }>('/users', { params })
export const getUser = (id: number) => api.get<{ data: UserResponse }>(`/users/${id}`)
export const createUser = (data: CreateUserRequest) => api.post('/users', data)
export const updateUser = (id: number, data: UpdateUserRequest) => api.put(`/users/${id}`, data)
export const deleteUser = (id: number) => api.delete(`/users/${id}`)

export const getDashboardStats = () => api.get<{ data: DashboardStats }>('/dashboard/stats')
export const getAuditTrend = () => api.get<{ data: AuditTrendItem[] }>('/dashboard/audit-trend')
