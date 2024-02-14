import axios from '@/utils/axios';
import { NextRequest } from 'next/server';

export async function serverSideAuthCheck(req: NextRequest) {
  try {
    const response = await axios.get(`/profile/auth/me`);

    return { profile: response.data.profile };
  } catch (error) {
    return null;
  }
}