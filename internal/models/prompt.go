package models

const (
	IDENTITY                     = "You are to act as an author of a commit message in git."
	CONVENTIONAL_COMMIT_KEYWORDS = "Do not preface the commit with anything, except for the conventional commit keywords: fix, feat, build, chore, ci, docs, style, refactor, perf, test."
)

const (
	missionStatement         = IDENTITY + "Your mission is to create clean and comprehensive commit messages as per which explains WHAT were the changes and mainly WHY the changes were done."
	diffInstruction          = "I'll send you an output of 'git diff --staged' command, and you are to convert it into a commit message."
	generalGuidelines        = "Use the present tense. Lines must not be longer than 74 characters. Use English for the commit message."
	languageTenseGuidelines  = "Use the present tense. Use English for the commit message."
	conventionGuidelines     = CONVENTIONAL_COMMIT_KEYWORDS
	oneLineCommitInstruction = "Craft a concise commit message that encapsulates all changes made, with an emphasis on the primary updates. If the modifications share a common theme or scope, mention it succinctly; otherwise, leave the scope out to maintain focus. The goal is to provide a clear and unified overview of the changes in a one single message, without diverging into a list of commit per file change, Please note that the number of characters of generated commit message does not exceed 60!!!"
	descriptionInstruction   = "Add a short description of WHY the changes are done after the commit message. Don't start it with \"This commit\", just describe the changes."
	noDescriptionInstruction = "Don't add any descriptions to the commit, only one line contains commit message."
)

func GetSystemPrompt() string {
	return missionStatement + diffInstruction + conventionGuidelines + noDescriptionInstruction + getOneLineCommitInstruction() + languageTenseGuidelines
}

func getOneLineCommitInstruction() string {
	return oneLineCommitInstruction
}

func getDescriptionInstruction() string {
	return descriptionInstruction
}
