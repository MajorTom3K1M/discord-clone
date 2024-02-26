import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function decodeJwtPayload(token: string): any {
  // Split the JWT string into parts
  const parts = token.split('.');
  if (parts.length !== 3) {
    throw new Error('The token is invalid');
  }

  // Decode the payload
  const payload = parts[1];
  const decodedPayload = Buffer.from(payload, 'base64url').toString('utf8');

  // Parse the JSON
  try {
    const parsedPayload = JSON.parse(decodedPayload);
    return parsedPayload;
  } catch (e) {
    throw new Error('Failed to parse the payload of the token');
  }
}