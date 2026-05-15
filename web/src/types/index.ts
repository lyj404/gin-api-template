export interface ResponseResult<T = any> {
  code: number
  message: string
  data?: T
}

export interface PaginationRequest {
  page?: number
  page_size?: number
  keyword?: string
  operator_id?: string
  order_by?: string
  sort?: 'asc' | 'desc'
}

export interface PaginationResponse<T = any> {
  page: number
  page_size: number
  total: number
  total_page: number
  data: T[]
}

export interface LoginRequest {
  email: string
  password: string
  captcha: string
}

export interface LoginResponse {
  accessToken: string
  refreshToken: string
}

export interface SignupRequest {
  name: string
  email: string
  password: string
}

export interface SignupResponse {
  accessToken: string
  refreshToken: string
}

export interface RefreshTokenRequest {
  refreshToken: string
}

export interface RefreshTokenResponse {
  accessToken: string
  refreshToken: string
}

export interface RoleRequest {
  name: string
  description?: string
}

export interface RoleResponse {
  id: number
  name: string
  description: string
  is_system: boolean
}

export interface ResourceBrief {
  id: number
  name: string
  type: string
  pattern: string
  method: string
  entity: string
  action: string
  description: string
}

export interface MenuResourceResponse {
  id: number
  menu_id: number
  resource_id: number
  resource?: ResourceBrief
}

export interface RoleResourceResponse {
  id: number
  role_id: number
  resource_id: number
  is_read: boolean
  is_write: boolean
  resource?: ResourceBrief
}

export interface RoleMenuResponse {
  id: number
  role_id: number
  menu_id: number
  menu?: MenuBrief
}

export interface MenuBrief {
  id: number
  name: string
  path: string
  icon: string
}

export interface RoleDetailResponse {
  id: number
  name: string
  description: string
  is_system: boolean
  resources?: RoleResourceResponse[]
  menus?: RoleMenuResponse[]
}

export interface CreateMenuRequest {
  name: string
  parent_id?: number | null
  path?: string
  icon?: string
  order_num?: number
  is_visible?: boolean
}

export interface UpdateMenuRequest {
  name?: string
  parent_id?: number | null
  path?: string
  icon?: string
  order_num?: number
  is_visible?: boolean
  status?: string
}

export interface MenuResponse {
  id: number
  name: string
  parent_id: number | null
  path: string
  icon: string
  order_num: number
  is_visible: boolean
  status: string
  resources?: ResourceBrief[]
  children?: MenuResponse[]
}

export interface MenuTreeNode {
  id: number
  name: string
  path: string
  icon: string
  order_num: number
  is_visible: boolean
  children?: MenuTreeNode[]
}

export interface CreateOrgUnitRequest {
  name: string
  parent_id?: number | null
}

export interface UpdateOrgUnitRequest {
  name?: string
  parent_id?: number | null
}

export interface OrgUnitResponse {
  id: number
  name: string
  parent_id: number | null
  path: string
  level: number
}

export interface ResourceResponse {
  id: number
  name: string
  type: string
  pattern: string
  method: string
  entity: string
  action: string
  description: string
}

export interface BindResourceRequest {
  is_write?: boolean
}

export interface BindMenuRequest {
  menu_id: number
}

export interface BindMenuResourceRequest {
  resource_id: number
}

export interface Permission {
  resource: string
  method: string
}

export interface UserPermissions {
  permissions: Permission[]
  org_scope: number[]
  resources: ResourceResponse[]
}

export interface UserResponse {
  id: number
  name: string
  email: string
  role_ids: number[]
  roles: string[]
}

export interface CreateUserRequest {
  name: string
  email: string
  password: string
  role_ids?: number[]
  org_unit_id?: number
}

export interface UpdateUserRequest {
  name?: string
  email?: string
  password?: string
  role_ids?: number[]
  org_unit_id?: number
}

export interface ProfileResponse {
  id: number
  name: string
  email: string
  created_at: string
  updated_at: string
}

export interface UpdateProfileRequest {
  name: string
  email: string
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export interface DashboardStats {
  user_count: number
  role_count: number
  menu_count: number
  audit_log_count: number
}

export interface AuditTrendItem {
  date: string
  count: number
}
