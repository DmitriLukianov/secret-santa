import { getHeaders, handleResponse } from './http.jsx';

/** Каталог/черновики товаров (страница wishlist-add). Путь без /v1 — как в текущем fetch на localhost:8080 */
const PRODUCTS_BASE = '/api';

export const createProduct = async (productData) => {
  const response = await fetch(`${PRODUCTS_BASE}/products`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(productData),
  });
  return handleResponse(response);
};

export const fetchProducts = async () => {
  const response = await fetch(`${PRODUCTS_BASE}/products`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};
