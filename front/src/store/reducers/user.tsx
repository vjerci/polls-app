import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { AppState } from "@/store/store";
import api from "@/lib/api";
import Router from "next/router";

export interface UserState {
	token: string;
	displayName: string;
	loginState: "loading" | "idle" | "failed";
	loginError: string;
	registerState: "loading" | "idle" | "failed";
	registerError: string;
	googleLoginState: "loading" | "idle" | "failed";
	googleLoginError: string;
}

const initialState: UserState = {
	token: "",
	displayName: "",
	loginState: "idle",
	loginError: "",
	registerState: "idle",
	registerError: "",
	googleLoginState: "idle",
	googleLoginError: "",
};

export const doGoogleLogin = createAsyncThunk(
	"user/doGoogleLogin",
	async (token: string, { rejectWithValue, dispatch }) => {
		const resp = await api.POST({
			path: "/auth/google/login",
			body: { token: token },
		});

		if (typeof resp.error === "undefined") {
			return dispatch(
				login({ token: resp.data.token, displayName: resp.data.name })
			);
		}

		return rejectWithValue(resp.error);
	}
);

export const doLogin = createAsyncThunk(
	"user/doLogin",
	async (userID: string, { rejectWithValue, dispatch }) => {
		const resp = await api.POST({
			path: "/auth/login",
			body: { user_id: userID },
		});

		if (typeof resp.error === "undefined") {
			return dispatch(
				login({ token: resp.data.token, displayName: resp.data.name })
			);
		}

		return rejectWithValue(resp.error);
	}
);

export const doRegister = createAsyncThunk(
	"user/doRegister",
	async (
		{ userID, displayName }: { userID: string; displayName: string },
		{ rejectWithValue, dispatch }
	) => {
		const resp = await api.PUT({
			path: "/auth/register",
			body: { user_id: userID, name: displayName },
		});

		if (typeof resp.error === "undefined") {
			return dispatch(login({ token: resp.data, displayName }));
		}

		return rejectWithValue(resp.error);
	}
);

export const userSlice = createSlice({
	name: "user",
	initialState,
	reducers: {
		login(
			state,
			action: { payload: { token: string; displayName: string } }
		) {
			state.token = action.payload.token;
			state.displayName = action.payload.displayName;

			Router.push("/polls");
		},
		logout(state) {
			state.token = "";
			state.displayName = "";

			Router.push("/");
		},
		clearLoginError(state) {
			state.loginState = "idle";
			state.loginError = "";
		},
		clearRegisterError(state) {
			state.registerState = "idle";
			state.registerError = "";
		},
	},

	extraReducers: (builder) => {
		builder
			.addCase(doLogin.pending, (state) => {
				state.loginState = "loading";
			})
			.addCase(doLogin.fulfilled, (state) => {
				state.loginState = "idle";
			})
			.addCase(doLogin.rejected, (state, action: any) => {
				state.loginState = "failed";
				state.loginError = action.payload;
			})

			.addCase(doRegister.pending, (state) => {
				state.registerState = "loading";
			})
			.addCase(doRegister.fulfilled, (state) => {
				state.registerState = "idle";
			})
			.addCase(doRegister.rejected, (state, action: any) => {
				state.registerState = "failed";
				state.registerError = action.payload;
			})

			.addCase(doGoogleLogin.pending, (state) => {
				state.googleLoginState = "loading";
			})
			.addCase(doGoogleLogin.fulfilled, (state) => {
				state.googleLoginState = "idle";
			})
			.addCase(doGoogleLogin.rejected, (state, action: any) => {
				state.googleLoginState = "failed";
				state.googleLoginError = action.payload;
			});
	},
});

export const { login, logout, clearLoginError, clearRegisterError } =
	userSlice.actions;

export const selectUserName = (state: AppState) => state.user.displayName;

export const selectLoginIsLoading = (state: AppState) =>
	state.user.loginState === "loading";
export const selectLoginError = (state: AppState) => state.user.loginError;

export const selectRegisterIsLoading = (state: AppState) =>
	state.user.registerState === "loading";
export const selectRegisterError = (state: AppState) =>
	state.user.registerError;

export const selectIsLoggedIn = (state: AppState) => state.user.token !== "";

export const selectGoogleLoginIsLoading = (state: AppState) =>
	state.user.googleLoginState === "loading";

export default userSlice.reducer;
