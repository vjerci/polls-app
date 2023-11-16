import api, { buildAuthHeaders } from "@/lib/api";
import { AppState } from "@/store/store";
import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";

export interface PollDetailsState {
	name: string;
	id: string;
	answers: Array<Answer>;
	loadingState: "idle" | "loading" | "failed";
	loadingError: string;
}

export interface Answer {
	name: string;
	id: string;
	votesCount: number;
	hasUserVoted: boolean;
}

const initialState: PollDetailsState = {
	name: "",
	id: "",
	answers: [],
	loadingState: "idle",
	loadingError: "",
};

export const fetchPollDetails = createAsyncThunk(
	"pollDetails/fetchPollDetails",
	async (
		{ pollId }: { pollId: string },
		{
			rejectWithValue,
			dispatch,
			getState,
		}: { rejectWithValue: any; dispatch: any; getState: any }
	) => {
		const state: AppState = getState();

		const resp = await api.GET({
			path: `/poll/${pollId}`,
			headers: buildAuthHeaders(state.user.token),
		});

		if (typeof resp.error === "undefined") {
			return dispatch(addPollDetails(resp.data));
		}

		return rejectWithValue(resp.error);
	}
);

export const pollDetailsVote = createAsyncThunk(
	"pollDetails/votePollDetails",
	async (
		{ answerId }: { answerId: string },
		{
			rejectWithValue,
			dispatch,
			getState,
		}: { rejectWithValue: any; dispatch: any; getState: any }
	) => {
		console.log("voting");

		const state: AppState = getState();
		const pollId = state.pollDetails.id;

		const resp = await api.POST({
			path: `/poll/${pollId}/vote`,
			headers: buildAuthHeaders(state.user.token),
			body: { answer_id: answerId },
		});

		if (typeof resp.error === "undefined") {
			return dispatch(addPollVote({ answerId: answerId }));
		}

		return rejectWithValue(resp.error);
	}
);

export const pollDetailsSlice = createSlice({
	name: "pollDetails",
	initialState,
	reducers: {
		addBasicPollInfo(state, action: { payload: { name: string } }) {
			state.answers = [];
			state.name = action.payload.name;
		},
		addPollDetails(
			state,
			action: {
				payload: {
					id: string;
					name: string;
					answers: Array<{
						name: string;
						id: string;
						votes_count: number;
					}>;
					user_vote: string;
				};
			}
		) {
			const answers: Array<Answer> = action.payload.answers.map((a) => ({
				name: a.name,
				votesCount: a.votes_count,
				id: a.id,
				hasUserVoted: a.id === action.payload.user_vote,
			}));

			state.id = action.payload.id;
			state.name = action.payload.name;
			state.answers = answers;
		},
		addPollVote(state, action: { payload: { answerId: string } }) {
			const previousAnswer = state.answers.find((a) => a.hasUserVoted);
			if (previousAnswer) {
				previousAnswer.hasUserVoted = false;
				previousAnswer.votesCount--;
			}

			const newAnswer = state.answers.find(
				(a) => a.id === action.payload.answerId
			);
			if (newAnswer) {
				newAnswer.hasUserVoted = true;
				newAnswer.votesCount++;
			}
		},
	},

	extraReducers: (builder) => {
		builder
			.addCase(fetchPollDetails.pending, (state) => {
				state.loadingState = "loading";
			})
			.addCase(fetchPollDetails.fulfilled, (state) => {
				state.loadingState = "idle";
			})
			.addCase(fetchPollDetails.rejected, (state, action: any) => {
				state.loadingState = "failed";
				state.loadingError = action.payload;
			})
			.addCase(pollDetailsVote.pending, (state) => {
				state.loadingState = "loading";
			})
			.addCase(pollDetailsVote.fulfilled, (state) => {
				state.loadingState = "idle";
			})
			.addCase(pollDetailsVote.rejected, (state, action: any) => {
				state.loadingState = "failed";
				state.loadingError = action.payload;
			});
	},
});

export const selectPollDetailsIsLoading = (state: AppState) =>
	state.pollDetails.loadingState == "loading";

export const selectPollDetailsName = (state: AppState) =>
	state.pollDetails.name;

export const selectPollDetailsAnswers = (state: AppState) =>
	state.pollDetails.answers;

export const { addPollDetails, addBasicPollInfo, addPollVote } =
	pollDetailsSlice.actions;
