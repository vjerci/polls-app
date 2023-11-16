import { configureStore, combineReducers } from "@reduxjs/toolkit";
import { userSlice } from "@/store/reducers/user";
import { pollsListSlice } from "@/store/reducers/polls_list";
import { pollDetailsSlice } from "@/store/reducers/poll_details";
import { createWrapper } from "next-redux-wrapper";
import { persistReducer, persistStore } from "redux-persist";
import storage from "redux-persist/lib/storage";

const rootReducer = combineReducers({
	[userSlice.name]: userSlice.reducer,
	[pollsListSlice.name]: pollsListSlice.reducer,
	[pollDetailsSlice.name]: pollDetailsSlice.reducer,
});

const makeConfiguredStore = () =>
	configureStore({
		reducer: rootReducer,
		devTools: true,
		middleware: (getDefaultMiddleware) =>
			getDefaultMiddleware({ serializableCheck: false }),
	});

export const makeStore = () => {
	const isServer = typeof window === "undefined";
	if (isServer) {
		return makeConfiguredStore();
	} else {
		const persistConfig = {
			key: "polls_app",
			whitelist: ["user"],
			storage,
		};
		const persistedReducer = persistReducer(persistConfig, rootReducer);
		let store: any = configureStore({
			reducer: persistedReducer,
			devTools: process.env.NODE_ENV !== "production",
			middleware: (getDefaultMiddleware) =>
				getDefaultMiddleware({ serializableCheck: false }),
		});
		store.__persistor = persistStore(store); // Nasty hack
		return store;
	}
};

export type AppStore = ReturnType<typeof makeConfiguredStore>;
export type AppState = ReturnType<AppStore["getState"]>;
export type AppDispatch = AppStore["dispatch"];

export const wrapper = createWrapper<AppStore>(makeStore);
