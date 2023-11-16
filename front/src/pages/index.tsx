"use client";

import { Inter } from "next/font/google";
import GoogleLogin from "@/containers/auth/google_login";
import Button from "@/components/common/button";
import { useRouter } from "next/router";

const inter = Inter({ subsets: ["latin"] });

export default function Home() {
	const router = useRouter();
	return (
		<div className={inter.className}>
			<h1 className="text-xl text-stone-100">Welcome to Polls app</h1>

			<h2 className="mt-2 text-l text-stone-100">
				In polls app you can create, browse and vote on different polls.
			</h2>

			<div className="mt-12 text-stone-200">
				To continue please use some of the actions below.
			</div>

			<nav className="my-8 max-w-xs">
				<h2 className="mb-4 text-l text-stone-100">
					With credentials:
				</h2>

				<Button
					name="Register"
					className="block w-full"
					onClick={() => {
						router.push("/register");
					}}
				/>
				<Button
					name="Login"
					className="block w-full"
					onClick={() => {
						router.push("/login");
					}}
				/>
			</nav>

			<GoogleLogin />
		</div>
	);
}
