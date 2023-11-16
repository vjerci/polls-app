"use client";

import Breadcrumb from "@/components/common/breadcrumb";
import Input from "@/components/common/input";
import Form from "@/components/common/form";
import { FormEvent, useEffect, useState } from "react";
import {
	clearLoginError,
	doLogin,
	selectLoginError,
	selectLoginIsLoading,
} from "@/store/reducers/user";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import GoogleLogin from "@/components/auth/google_login_disabled";
import GoogleLoginDisabled from "@/components/auth/google_login_disabled";

export default function Page() {
	const dispatch = useAppDispatch();
	const loginError = useAppSelector(selectLoginError);
	const loading = useAppSelector(selectLoginIsLoading);

	useEffect(() => {
		dispatch(clearLoginError());
	}, [dispatch]);

	const onSubmit = (e: FormEvent) => {
		e.preventDefault();
		dispatch(doLogin(userID));
	};

	const [userID, setUserID] = useState("");
	const handleUserIDChange = (e: any) => {
		setUserID(e.target.value);
		dispatch(clearLoginError());
	};

	return (
		<div>
			<Breadcrumb homeHref="/" links={[{ name: "Login" }]} />
			<h1 className="text-xl my-8 text-stone-100">Login</h1>
			<GoogleLoginDisabled />

			<Form submitName="Login" onSubmit={onSubmit} loading={loading}>
				<Input
					placeholder="User ID"
					value={userID}
					error={loginError}
					onChange={handleUserIDChange}
				/>
			</Form>
		</div>
	);
}
