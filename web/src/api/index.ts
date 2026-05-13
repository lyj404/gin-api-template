import api from '@/utils/http'
import type {
  LoginRequest, LoginResponse,
  SignupRequest, SignupResponse,
  RefreshTokenRequest, RefreshTokenResponse,
  RoleRequest, RoleResponse,
  PaginationRequest, PaginationResponse,
  CreateMenuRequest, UpdateMenuRequest, MenuResponse, MenuTreeNode,
  CreateOrgUnitRequest, UpdateOrgUnitRequest, OrgUnitResponse,
  ResourceResponse,
  BindResourceRequest,
  UserPermissions,
  UserResponse, CreateUserRequest, UpdateUserRequest,
  DashboardStats, AuditTrendItem
} from '@/types'

export const login = (data: LoginRequest) => api.post('/login', data)
export const signup = (data: SignupRequest) => api.post('/signup', data)
export const refreshToken = (data: RefreshTokenRequest) => api.post('/refresh-token', data)
export const getCaptcha = () => api.get('/captcha', { responseType: 'blob' })

export const getRoles = (params?: PaginationRequest) => api.get<PaginationResponse<RoleResponse>>('/roles', { params })
export const getRole = (id: number) => api.get<{ data: RoleResponse }>(`/roles/${id}`)
export const createRole = (data: RoleRequest) => api.post<{ data: RoleResponse }>('/roles', data)
export const updateRole = (id: number, data: RoleRequest) => api.put(`/roles/${id}`, data)
export const deleteRole = (id: number) => api.delete(`/roles/${id}`)

export const getMenuTree = () => api.get<{ data: MenuTreeNode[] }>('/menus/tree')
export const getMenu = (id: number) => api.get<{ data: MenuResponse }>(`/menus/${id}`)
export const createMenu = (data: CreateMenuRequest) => api.post('/menus', data)
export const updateMenu = (id: number, data: UpdateMenuRequest) => api.put(`/menus/${id}`, data)
export const deleteMenu = (id: number) => api.delete(`/menus/${id}`)

export const getOrgUnits = (params?: PaginationRequest) => api.get<PaginationResponse<OrgUnitResponse>>('/org-units', { params })
export const getOrgTree = () => api.get<{ data: OrgUnitResponse[] }>('/org-units/tree')
export const getOrgUnit = (id: number) => api.get<{ data: OrgUnitResponse }>(`/org-units/${id}`)
export const createOrgUnit = (data: CreateOrgUnitRequest) => api.post('/org-units', data)
export const updateOrgUnit = (id: number, data: UpdateOrgUnitRequest) => api.put(`/org-units/${id}`, data)
export const deleteOrgUnit = (id: number) => api.delete(`/org-units/${id}`)

export const getResources = (params?: PaginationRequest) => api.get<PaginationResponse<ResourceResponse>>('/resources', { params })
export const bindResource = (roleId: number, resourceId: number, data: BindResourceRequest) => api.post(`/roles/${roleId}/resources`, { ...data, resource_id: resourceId })

export const getAuditLogs = (params?: PaginationRequest) => api.get<PaginationResponse<any>>('/audit-logs', { params })
export const getAuditLogsByTarget = (params: { target_type: string; target_id: string }) => api.get<{ data: any[] }>('/audit-logs/target', { params })
export const getAuditLogsByTime = (params: { start_time: string; end_time: string }) => api.get<{ data: any[] }>('/audit-logs/time', { params })

export const getUserPermissions = () => api.get<{ data: UserPermissions }>('/user/permissions')
export const getUserMenus = () => api.get<{ data: { menus: MenuTreeNode[] } }>('/user/menus')

export const getUsers = (params?: PaginationRequest) => api.get<PaginationResponse<UserResponse>>('/users', { params })
export const getUser = (id: number) => api.get<{ data: UserResponse }>(`/users/${id}`)
export const createUser = (data: CreateUserRequest) => api.post('/users', data)
export const updateUser = (id: number, data: UpdateUserRequest) => api.put(`/users/${id}`, data)
export const deleteUser = (id: number) => api.delete(`/users/${id}`)

export const getDashboardStats = () => api.get<{ data: DashboardStats }>('/dashboard/stats')
export const getAuditTrend = () => api.get<{ data: AuditTrendItem[] }>('/dashboard/audit-trend')
