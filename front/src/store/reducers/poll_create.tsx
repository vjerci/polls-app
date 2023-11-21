import api, { buildAuthHeaders } from "@/lib/api";
import { AppState } from "@/store/store";
import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import Router from "next/router";

export interface PollCreateState {
	loadingState: "idle" | "loading" | "failed";
	loadingError: string;
}

const initialState: PollCreateState = {
	loadingState: "idle",
	loadingError: "",
};

export const createPoll = createAsyncThunk(
	"pollCreate/createPoll",
	async (
		{ pollName, answers }: { pollName: string; answers: Array<string> },
		{
			rejectWithValue,
			getState,
		}: { rejectWithValue: any; dispatch: any; getState: any }
	) => {
		const state: AppState = getState();

		const resp = await api.PUT({
			path: `/poll`,
			headers: buildAuthHeaders(state.user.token),
			body: { name: pollName, answers },
		});

		if (typeof resp.error === "undefined") {
			return Router.push(`/polls/${resp.data.id}`);
		}

		return rejectWithValue(resp.error);
	}
);

export const pollCreateSlice = createSlice({
	name: "pollCreate",
	initialState,
	reducers: {},

	extraReducers: (builder) => {
		builder
			.addCase(createPoll.pending, (state) => {
				state.loadingState = "loading";
			})
			.addCase(createPoll.fulfilled, (state) => {
				state.loadingState = "idle";
			})
			.addCase(createPoll.rejected, (state, action: any) => {
				state.loadingState = "failed";
				state.loadingError = action.payload;
			});
	},
});


export const selectPollCreateIsLoading = (state: AppState) =>
	state.pollDetails.loadingState == "loading";