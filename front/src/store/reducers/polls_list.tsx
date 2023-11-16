import api, { buildAuthHeaders } from "@/lib/api";
import { AppState } from "@/store/store";
import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import App from "next/app";

export interface PollsListState {
	polls: Array<Poll>;
	page: number;
	hasNextPage: boolean;
	loadingState: "idle" | "loading" | "failed";
	loadingError: string;
}

export interface Poll {
	name: string;
	id: string;
}

const initialState: PollsListState = {
	polls: [],
	page: 0,
	hasNextPage: false,
	loadingState: "idle",
	loadingError: "",
};

export const fetchPollsList = createAsyncThunk(
	"pollsList/fetchData",
	async (
		{ page }: { page: number },
		{
			rejectWithValue,
			dispatch,
			getState,
		}: { rejectWithValue: any; dispatch: any; getState: any }
	) => {
		const state: AppState = getState();

		const resp = await api.GET({
			path: `/poll?page=${page}`,
			headers: buildAuthHeaders(state.user.token),
		});

		if (typeof resp.error === "undefined") {
			return dispatch(
				addPolls({
					page: page,
					polls: resp.data.polls,
					hasNext: resp.data.has_next,
				})
			);
		}

		return rejectWithValue(resp.error);
	}
);

export const pollsListSlice = createSlice({
	name: "pollsList",
	initialState,
	reducers: {
		addPolls(
			state,
			action: {
				payload: { page: number; polls: Array<Poll>; hasNext: boolean };
			}
		) {
			state.page = action.payload.page;
			state.polls = action.payload.polls;
			state.hasNextPage = action.payload.hasNext;
		},
	},

	extraReducers: (builder) => {
		builder
			.addCase(fetchPollsList.pending, (state) => {
				state.loadingState = "loading";
			})
			.addCase(fetchPollsList.fulfilled, (state) => {
				state.loadingState = "idle";
			})
			.addCase(fetchPollsList.rejected, (state, action: any) => {
				state.loadingState = "failed";
				state.loadingError = action.payload;
			});
	},
});

export const selectPollsListIsLoading = (state: AppState) =>
	state.pollsList.loadingState == "loading";
export const selectPollsListHasNextPage = (state: AppState) =>
	state.pollsList.hasNextPage;
export const selectPollsListData = (state: AppState) => state.pollsList.polls;
export const selectPollsListError = (state: AppState) =>
	state.pollsList.loadingError;
export const selectPollListPage = (state: AppState) => state.pollsList.page;

export const { addPolls } = pollsListSlice.actions;
