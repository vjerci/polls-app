import Button from "@/components/common/button";
import Error from "@/components/common/error";
import Form from "@/components/common/form";
import Input from "@/components/common/input";
import Link from "@/components/common/link";
import LogoutButton from "@/containers/polls/logout_button";
import PollBreadCrumb from "@/containers/polls/poll_breadcrumb";
import { useAppDispatch } from "@/store/hooks";
import { createPoll, selectPollCreateIsLoading } from "@/store/reducers/poll_create";
import { useState } from "react";
import { useSelector } from "react-redux";

interface Answer {
	name: string;
	value: string;
	error: string;
}

export default function Page() {
	const dispatch = useAppDispatch();

	const loading = useSelector(selectPollCreateIsLoading);

	const [answers, setAnswers] = useState([
		{ name: "Answer 1", value: "", error: "" },
		{ name: "Answer 2", value: "", error: "" },
	]);

	const [questionError, setQuestionError] = useState("");
	const [question, setPollQuestion] = useState("");
	const handlePollQuestion = (e: any) => {
		setPollQuestion(e.target.value);
		setQuestionError("");
	};

	const modifyAnswer = (index: number, newValue: string) => {
		setAnswers(
			answers.map((answer: Answer, localIndex: number): Answer => {
				if (index === localIndex) {
					return { ...answer, value: newValue, error: "" };
				}

				return answer;
			})
		);
	};

	const removeAnswer = (index: number) => {
		setAnswers(
			answers.filter((answer, localIndex) => localIndex !== index)
		);
	};

	const addAnswer = () => {
		const index = answers.length + 1;
		setAnswers([
			...answers,
			{ name: `Answer ${index}`, value: "", error: "" },
		]);
		setFormError("");
	};

	const [formError, setFormError] = useState("");
	const initPollCreation = () => {
		let hasError = false;

		if (answers.length < 2) {
			hasError = true;
			setFormError("Minimum of 2 answers is required");
		}

		setAnswers(
			answers.map((answer) => {
				if (answer.value.length === 0) {
					hasError = true;
					return { ...answer, error: "Answer can't be empty" };
				}
				return answer;
			})
		);

		if (question.length === 0) {
			hasError = true;
			setQuestionError("Question can't be empty");
		}

		if (!hasError) {
			dispatch(
				createPoll({
					pollName: question,
					answers: answers.map((a) => a.value),
				})
			);
		}
	};

	return (
		<>
			<nav className="mb-4 grid">
				<LogoutButton />
			</nav>
			<PollBreadCrumb currentName="Poll Creation" />

			<h2 className="text-2xl text-slate-100 pb-4">Create Poll</h2>

			<p className="text-slate-500">To create a valid pool:</p>
			<ol className="my-2 text-slate-500">
				<li className="ml-2"> - Have non empty question</li>
				<li className="ml-2">- Provide at least 2 non empty answers</li>
			</ol>

			<Form
				onSubmit={(e) => {
					e.preventDefault();
					initPollCreation();
				}}
				submitName="Create poll"
				loading={loading}
			>
				<Input
					placeholder="Poll question"
					onChange={handlePollQuestion}
					value={question}
					error={questionError}
				/>
				<hr className="py-2" />
				{answers.map((answer, index: number) => (
					<div key={index}>
						<Input
							placeholder={answer.name}
							onChange={(e: any) =>
								modifyAnswer(index, e.target.value)
							}
							value={answer.value}
							error={answer.error}
						/>

						<Button
							onClick={() => removeAnswer(index)}
							name="Remove answer"
						></Button>

						<hr className="py-2" />
					</div>
				))}

				<Button onClick={() => addAnswer()} name="Add answer"></Button>

				<Error error={formError} />
			</Form>
		</>
	);
}
