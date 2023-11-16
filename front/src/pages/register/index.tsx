"use client";

import GoogleLoginDisabled from "@/components/auth/google_login_disabled";
import GoogleLogin from "@/components/auth/google_login_disabled";
import Breadcrumb from "@/components/common/breadcrumb";
import Error from "@/components/common/error";
import Form from "@/components/common/form";
import Input from "@/components/common/input";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import {
	clearRegisterError,
	doRegister,
	selectRegisterError,
	selectRegisterIsLoading,
} from "@/store/reducers/user";
import { FormEvent, useEffect, useState } from "react";

export default function Page() {
	const dispatch = useAppDispatch();

	const loading = useAppSelector(selectRegisterIsLoading);
	const registerError = useAppSelector(selectRegisterError);

	useEffect(() => {
		dispatch(clearRegisterError());
	}, [dispatch]);

	const onSubmit = (e: FormEvent) => {
		e.preventDefault();
		dispatch(doRegister({ userID: userID, displayName: displayName }));
	};

	const [userID, setUserID] = useState("");
	const handleUserIDChange = (e: any) => {
		setUserID(e.target.value);
		dispatch(clearRegisterError());
	};

	const [displayName, setDisplayName] = useState("");
	const handleDisplayNameChange = (e: any) => {
		setDisplayName(e.target.value);
		dispatch(clearRegisterError());
	};

	return (
		<div>
			<Breadcrumb homeHref="/" links={[{ name: "Register" }]} />

			<h1 className="text-xl my-8 text-stone-100">Register</h1>

			<GoogleLoginDisabled />

			<Form submitName="Register" onSubmit={onSubmit} loading={loading}>
				<Input
					placeholder="User ID"
					onChange={handleUserIDChange}
					value={userID}
				/>
				<Input
					placeholder="Display Name"
					onChange={handleDisplayNameChange}
					value={displayName}
				/>

				<Error error={registerError} />
			</Form>
		</div>
	);
}
