import { http } from '../http';

// 产品列表查询参数
export interface ProductQueryParams {
  page?: number;
  pageSize?: number;
  categoryID?: string;
  name?: string;

}

// 产品创建请求
export interface ProductCreateRequest {
  categoryID?: string;
  name: string;
  sort?: number;

}

// 产品更新请求
export interface ProductUpdateRequest {
  id: number;
  categoryID?: string;
  name?: string;
  sort?: number;

}

// 产品响应
export interface ProductResponse {
  id: number;
  createdAt: string;
  updatedAt: string;
  categoryID: string;
  name: string;
  sort: number;

}

// 产品列表响应
export interface ProductListResponse {
  total: number;
  list: ProductResponse[];
}

// 获取产品列表
export const getProducts = (params: ProductQueryParams) => {
  return http.request<ProductListResponse>({
    url: '/product/list',
    method: 'GET',
    params,
  });
};

// 获取产品详情
export const getProduct = (id: number) => {
  return http.request<ProductResponse>({
    url: '/product/detail/' + id,
    method: 'GET',
  });
};

// 创建产品
export const createProduct = (data: ProductCreateRequest) => {
  return http.request<void>({
    url: '/product/create',
    method: 'POST',
    data,
  });
};

// 更新产品
export const updateProduct = (data: ProductUpdateRequest) => {
  return http.request<void>({
    url: '/product/update',
    method: 'PUT',
    data,
  });
};

// 删除产品
export const deleteProduct = (id: number) => {
  return http.request<void>({
    url: '/product/delete/' + id,
    method: 'DELETE',
  });
};
