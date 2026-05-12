export interface ResponseResult<T = any> {
  code: number
  message: string
  data?: T
}

export interface PaginationRequest {
  page?: number
  page_size?: number
  keyword?: string
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

export interface CreateMenuRequest {
  name: string
  parent_id?: number | null
  path?: string
  component?: string
  icon?: string
  order_num?: number
  resource_id: number
  is_visible?: boolean
}

export interface UpdateMenuRequest {
  name?: string
  parent_id?: number | null
  path?: string
  component?: string
  icon?: string
  order_num?: number
  resource_id?: number
  is_visible?: boolean
  status?: string
}

export interface MenuResponse {
  id: number
  name: string
  parent_id: number | null
  path: string
  component: string
  icon: string
  order_num: number
  resource_id: number
  is_visible: boolean
  status: string
  children?: MenuResponse[]
}

export interface MenuTreeNode {
  id: number
  name: string
  path: string
  component: string
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
  role_id: number
  is_write?: boolean
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
