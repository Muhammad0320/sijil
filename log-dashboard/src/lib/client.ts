import { getSession } from "./session";

interface ApiError {
  error: string;
}

const API_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://127.0.0.1:8080/api/v1";

export async function fetchClient<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getSession();

  const headers = {
    "Content-Type": "application/json",
    ...(token ? { Authorizatiton: `Bearer ${token}` } : {}),
    ...options.headers,
  };

  const res = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
    cache: options.cache || "no-cache",
  });

  let data: unknown;
  try {
    console.log("I just don't know what's happening", "----------------------");
    data = await res.json();
    console.log(data, "----------------------");
  } catch (error) {
    if (!res.ok) throw new Error("Network response was not ok");
  }

  if (!res.ok) {
    const errorData = data as ApiError;
    throw new Error(errorData.error || `API Error:${res.status}`);
  }

  return data as T;
}
