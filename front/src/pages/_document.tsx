import { Html, Head, Main, NextScript } from "next/document";

export default function Document() {
	return (
		<Html lang="en">
			<Head>
				<meta
					name="viewport"
					content="width=device-width,initial-scale=1,minimum-scale=1"
				/>
			</Head>
			<body className="grid">
				<main className="min-h-screen justify-self-center w-auto p-4 lg:p-24 lg:w-[960px] ">
					<Main />
				</main>
				<NextScript />
			</body>
		</Html>
	);
}
