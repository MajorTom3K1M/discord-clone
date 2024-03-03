import axios from "@/utils/axios";

export const clientRefreshToken = async () => {
    try {
        const response = await axios.get("/refresh", {
            withCredentials: true
        });
        return response.data.accessToken;
    } catch (err) {
        throw err;
    }
};