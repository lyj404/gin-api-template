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
  id: string
  name: string
  description: string
  is_system: boolean
}

export interface ResourceBrief {
  id: string
  name: string
  type: string
  pattern: string
  method: string
  entity: string
  action: string
  description: string
}

export interface MenuResourceResponse {
  id: string
  menu_id: string
  resource_id: string
  resource?: ResourceBrief
}

export interface RoleResourceResponse {
  id: string
  role_id: string
  resource_id: string
  is_read: boolean
  is_write: boolean
  resource?: ResourceBrief
}

export interface RoleMenuResponse {
  id: string
  role_id: string
  menu_id: string
  menu?: MenuBrief
}

export interface MenuBrief {
  id: string
  name: string
  path: string
  icon: string
}

export interface RoleDetailResponse {
  id: string
  name: string
  description: string
  is_system: boolean
  resources?: RoleResourceResponse[]
  menus?: RoleMenuResponse[]
}

export interface CreateMenuRequest {
  name: string
  parent_id?: string | null
  path?: string
  icon?: string
  order_num?: number
  is_visible?: boolean
  status?: string
}

export interface UpdateMenuRequest {
  name?: string
  parent_id?: string | null
  path?: string
  icon?: string
  order_num?: number
  is_visible?: boolean
  status?: string
}

export interface MenuResponse {
  id: string
  name: string
  parent_id: string | null
  path: string
  icon: string
  order_num: number
  is_visible: boolean
  status: string
  resources?: ResourceBrief[]
  children?: MenuResponse[]
}

export interface MenuTreeNode {
  id: string
  name: string
  path: string
  icon: string
  order_num: number
  is_visible: boolean
  status?: string
  children?: MenuTreeNode[]
}

export interface CreateOrgUnitRequest {
  name: string
  parent_id?: string | null
}

export interface UpdateOrgUnitRequest {
  name?: string
  parent_id?: string | null
}

export interface OrgUnitResponse {
  id: string
  name: string
  parent_id: string | null
  path: string
  level: number
}

export interface ResourceResponse {
  id: string
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
  menu_id: string
}

export interface BindMenuResourceRequest {
  resource_id: string
}

export interface Permission {
  resource: string
  method: string
}

export interface UserPermissions {
  permissions: Permission[]
  org_scope: string[]
  resources: ResourceResponse[]
}

export interface UserResponse {
  id: string
  name: string
  email: string
  role_ids: string[]
  roles: string[]
}

export interface CreateUserRequest {
  name: string
  email: string
  password: string
  role_ids?: string[]
  org_unit_id?: string
}

export interface UpdateUserRequest {
  name?: string
  email?: string
  password?: string
  role_ids?: string[]
  org_unit_id?: string
}

export interface ProfileResponse {
  id: string
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

// 字典管理
export interface DictResponse {
	id: string
	name: string
	type: string
	status: number
	desc: string
	details?: DictDetailResponse[]
}

export interface DictDetailResponse {
	id: string
	dict_id: string
	label: string
	value: string
	sort: number
	status: number
	remark: string
}

export interface CreateDictRequest {
	name: string
	type: string
	status?: number
	desc?: string
}

export interface UpdateDictRequest {
	name?: string
	type?: string
	status?: number
	desc?: string
}

export interface CreateDictDetailRequest {
	dict_id: string
	label: string
	value: string
	sort?: number
	status?: number
	remark?: string
}

export interface UpdateDictDetailRequest {
	label?: string
	value?: string
	sort?: number
	status?: number
	remark?: string
}
