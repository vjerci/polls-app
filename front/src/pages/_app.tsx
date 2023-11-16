import "@/styles/globals.css";

import type { AppProps } from "next/app";
import { wrapper } from "@/store/store";
import { PersistGate } from "redux-persist/integration/react";
import { useStore } from "react-redux";
import RouteProtector from "@/containers/route_protector";

function App({ Component, pageProps }: AppProps) {
	const store: any = useStore();
	return (
		<PersistGate
			persistor={store.__persistor}
			loading={<div>Loading...</div>}
		>
			<RouteProtector>
				<Component {...pageProps} />
			</RouteProtector>
		</PersistGate>
	);
}
export default wrapper.withRedux(App);
