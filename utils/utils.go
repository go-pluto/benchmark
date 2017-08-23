package utils

import (
	"fmt"
	"strings"

	"math/rand"
)

// Variables

// msgHeaders offers five choices for email headers.
var msgHeaders [5]string = [5]string{
	fmt.Sprintf("Date: Mon, 7 Feb 1994 21:52:25 -0800 (PST)\r\nFrom: Fred Foobar <foobar@Blurdybloop.COM>\r\nSubject: Important read\r\nTo: mooch@owatagu.siam.edu\r\nMessage-Id: <B27397-0100000@Blurdybloop.COM>\r\nMIME-Version: 1.0\r\nContent-Type: TEXT/PLAIN; CHARSET=US-ASCII\r\n\r\n"),
	fmt.Sprintf("Date: Mon, 7 Feb 1998 21:52:25 -0800 (PST)\r\nFrom: Sandra Silvester <foobar@Blurdybloop.COM>\r\nSubject: Have you seen the UN's Universal Declaration of Human Rights?\r\nTo: mooch@owatagu.siam.edu\r\nMessage-Id: <B27397-0100000@Blurdybloop.COM>\r\nMIME-Version: 1.0\r\nContent-Type: TEXT/PLAIN; CHARSET=US-ASCII\r\n\r\n"),
	fmt.Sprintf("Date: Mon, 7 Feb 2002 21:52:25 -0800 (PST)\r\nFrom: Ernie Ernest <foobar@Blurdybloop.COM>\r\nSubject: just found this on the Internet\r\nTo: mooch@owatagu.siam.edu\r\nMessage-Id: <B27397-0100000@Blurdybloop.COM>\r\nMIME-Version: 1.0\r\nContent-Type: TEXT/PLAIN; CHARSET=US-ASCII\r\n\r\n"),
	fmt.Sprintf("Date: Mon, 7 Feb 2006 21:52:25 -0800 (PST)\r\nFrom: Maria Magnolia <foobar@Blurdybloop.COM>\r\nSubject: Please forward to all your friends\r\nTo: mooch@owatagu.siam.edu\r\nMessage-Id: <B27397-0100000@Blurdybloop.COM>\r\nMIME-Version: 1.0\r\nContent-Type: TEXT/PLAIN; CHARSET=US-ASCII\r\n\r\n"),
	fmt.Sprintf("Date: Mon, 7 Feb 2010 21:52:25 -0800 (PST)\r\nFrom: Clara Clarabella <foobar@Blurdybloop.COM>\r\nSubject: seen dis?\r\nTo: mooch@owatagu.siam.edu\r\nMessage-Id: <B27397-0100000@Blurdybloop.COM>\r\nMIME-Version: 1.0\r\nContent-Type: TEXT/PLAIN; CHARSET=US-ASCII\r\n\r\n"),
}

// msgBodyLines in total represents the
// "Universal Declaration of Human Rights" of the United Nations,
// see: http://www.un.org/en/universal-declaration-human-rights
var msgBodyLines [38]string = [38]string{
	"Whereas recognition of the inherent dignity and of the equal and inalienable rights of all members of the human family is the foundation of freedom, justice and peace in the world,\r\n",
	"Whereas disregard and contempt for human rights have resulted in barbarous acts which have outraged the conscience of mankind, and the advent of a world in which human beings shall enjoy freedom of speech and belief and freedom from fear and want has been proclaimed as the highest aspiration of the common people,\r\n",
	"Whereas it is essential, if man is not to be compelled to have recourse, as a last resort, to rebellion against tyranny and oppression, that human rights should be protected by the rule of law,\r\n",
	"Whereas it is essential to promote the development of friendly relations between nations,\r\n",
	"Whereas the peoples of the United Nations have in the Charter reaffirmed their faith in fundamental human rights, in the dignity and worth of the human person and in the equal rights of men and women and have determined to promote social progress and better standards of life in larger freedom,\r\n",
	"Whereas Member States have pledged themselves to achieve, in co-operation with the United Nations, the promotion of universal respect for and observance of human rights and fundamental freedoms,\r\n",
	"Whereas a common understanding of these rights and freedoms is of the greatest importance for the full realization of this pledge,\r\n",
	"Now, Therefore THE GENERAL ASSEMBLY proclaims THIS UNIVERSAL DECLARATION OF HUMAN RIGHTS as a common standard of achievement for all peoples and all nations, to the end that every individual and every organ of society, keeping this Declaration constantly in mind, shall strive by teaching and education to promote respect for these rights and freedoms and by progressive measures, national and international, to secure their universal and effective recognition and observance, both among the peoples of Member States themselves and among the peoples of territories under their jurisdiction.\r\n",
	"Article 1: All human beings are born free and equal in dignity and rights. They are endowed with reason and conscience and should act towards one another in a spirit of brotherhood.\r\n",
	"Article 2: Everyone is entitled to all the rights and freedoms set forth in this Declaration, without distinction of any kind, such as race, colour, sex, language, religion, political or other opinion, national or social origin, property, birth or other status. Furthermore, no distinction shall be made on the basis of the political, jurisdictional or international status of the country or territory to which a person belongs, whether it be independent, trust, non-self-governing or under any other limitation of sovereignty.\r\n",
	"Article 3: Everyone has the right to life, liberty and security of person.\r\n",
	"Article 4: No one shall be held in slavery or servitude; slavery and the slave trade shall be prohibited in all their forms.\r\n",
	"Article 5: No one shall be subjected to torture or to cruel, inhuman or degrading treatment or punishment.\r\n",
	"Article 6: Everyone has the right to recognition everywhere as a person before the law.\r\n",
	"Article 7: All are equal before the law and are entitled without any discrimination to equal protection of the law. All are entitled to equal protection against any discrimination in violation of this Declaration and against any incitement to such discrimination.\r\n",
	"Article 8: Everyone has the right to an effective remedy by the competent national tribunals for acts violating the fundamental rights granted him by the constitution or by law.\r\n",
	"Article 9: No one shall be subjected to arbitrary arrest, detention or exile.\r\n",
	"Article 10: Everyone is entitled in full equality to a fair and public hearing by an independent and impartial tribunal, in the determination of his rights and obligations and of any criminal charge against him.\r\n",
	"Article 11: (1) Everyone charged with a penal offence has the right to be presumed innocent until proved guilty according to law in a public trial at which he has had all the guarantees necessary for his defence. (2) No one shall be held guilty of any penal offence on account of any act or omission which did not constitute a penal offence, under national or international law, at the time when it was committed. Nor shall a heavier penalty be imposed than the one that was applicable at the time the penal offence was committed.\r\n",
	"Article 12: No one shall be subjected to arbitrary interference with his privacy, family, home or correspondence, nor to attacks upon his honour and reputation. Everyone has the right to the protection of the law against such interference or attacks.\r\n",
	"Article 13: (1) Everyone has the right to freedom of movement and residence within the borders of each state. (2) Everyone has the right to leave any country, including his own, and to return to his country.\r\n",
	"Article 14: (1) Everyone has the right to seek and to enjoy in other countries asylum from persecution. (2) This right may not be invoked in the case of prosecutions genuinely arising from non-political crimes or from acts contrary to the purposes and principles of the United Nations.\r\n",
	"Article 15: (1) Everyone has the right to a nationality. (2) No one shall be arbitrarily deprived of his nationality nor denied the right to change his nationality.\r\n",
	"Article 16: (1) Men and women of full age, without any limitation due to race, nationality or religion, have the right to marry and to found a family. They are entitled to equal rights as to marriage, during marriage and at its dissolution. (2) Marriage shall be entered into only with the free and full consent of the intending spouses. (3) The family is the natural and fundamental group unit of society and is entitled to protection by society and the State.\r\n",
	"Article 17: (1) Everyone has the right to own property alone as well as in association with others. (2) No one shall be arbitrarily deprived of his property.\r\n",
	"Article 18: Everyone has the right to freedom of thought, conscience and religion; this right includes freedom to change his religion or belief, and freedom, either alone or in community with others and in public or private, to manifest his religion or belief in teaching, practice, worship and observance.\r\n",
	"Article 19: Everyone has the right to freedom of opinion and expression; this right includes freedom to hold opinions without interference and to seek, receive and impart information and ideas through any media and regardless of frontiers.\r\n",
	"Article 20: (1) Everyone has the right to freedom of peaceful assembly and association. (2) No one may be compelled to belong to an association.\r\n",
	"Article 21: (1) Everyone has the right to take part in the government of his country, directly or through freely chosen representatives. (2) Everyone has the right of equal access to public service in his country. (3) The will of the people shall be the basis of the authority of government; this will shall be expressed in periodic and genuine elections which shall be by universal and equal suffrage and shall be held by secret vote or by equivalent free voting procedures.\r\n",
	"Article 22: Everyone, as a member of society, has the right to social security and is entitled to realization, through national effort and international co-operation and in accordance with the organization and resources of each State, of the economic, social and cultural rights indispensable for his dignity and the free development of his personality.\r\n",
	"Article 23: (1) Everyone has the right to work, to free choice of employment, to just and favourable conditions of work and to protection against unemployment. (2) Everyone, without any discrimination, has the right to equal pay for equal work. (3) Everyone who works has the right to just and favourable remuneration ensuring for himself and his family an existence worthy of human dignity, and supplemented, if necessary, by other means of social protection. (4) Everyone has the right to form and to join trade unions for the protection of his interests.\r\n",
	"Article 24: Everyone has the right to rest and leisure, including reasonable limitation of working hours and periodic holidays with pay.\r\n",
	"Article 25: (1) Everyone has the right to a standard of living adequate for the health and well-being of himself and of his family, including food, clothing, housing and medical care and necessary social services, and the right to security in the event of unemployment, sickness, disability, widowhood, old age or other lack of livelihood in circumstances beyond his control. (2) Motherhood and childhood are entitled to special care and assistance. All children, whether born in or out of wedlock, shall enjoy the same social protection.\r\n",
	"Article 26: (1) Everyone has the right to education. Education shall be free, at least in the elementary and fundamental stages. Elementary education shall be compulsory. Technical and professional education shall be made generally available and higher education shall be equally accessible to all on the basis of merit. (2) Education shall be directed to the full development of the human personality and to the strengthening of respect for human rights and fundamental freedoms. It shall promote understanding, tolerance and friendship among all nations, racial or religious groups, and shall further the activities of the United Nations for the maintenance of peace. (3) Parents have a prior right to choose the kind of education that shall be given to their children.\r\n",
	"Article 27: (1) Everyone has the right freely to participate in the cultural life of the community, to enjoy the arts and to share in scientific advancement and its benefits. (2) Everyone has the right to the protection of the moral and material interests resulting from any scientific, literary or artistic production of which he is the author.\r\n",
	"Article 28: Everyone is entitled to a social and international order in which the rights and freedoms set forth in this Declaration can be fully realized.\r\n",
	"Article 29: (1) Everyone has duties to the community in which alone the free and full development of his personality is possible. (2) In the exercise of his rights and freedoms, everyone shall be subject only to such limitations as are determined by law solely for the purpose of securing due recognition and respect for the rights and freedoms of others and of meeting the just requirements of morality, public order and the general welfare in a democratic society. (3) These rights and freedoms may in no case be exercised contrary to the purposes and principles of the United Nations.\r\n",
	"Article 30: Nothing in this Declaration may be interpreted as implying for any State, group or person any right to engage in any activity or to perform any act aimed at the destruction of any of the rights and freedoms set forth herein.\r\n",
}

// Functions

// GenerateString returns a random string from the
// alphabet [a-z,0-9] of length "strlen".
func GenerateString(strlen int) string {

	// Define alphabet.
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

	result := ""
	for i := 0; i < strlen; i++ {
		index := rand.Intn(len(chars))
		result += chars[index:(index + 1)]
	}

	return result
}

// GenerateFlag returns a random choice of message flags.
func GenerateFlags() (string, []string) {

	// Define alphabet.
	flags := []string{"\\Seen", "\\Answered", "\\Flagged", "\\Deleted", "\\Draft"}

	numFlags := rand.Intn(len(flags)) + 1

	// Generate an array of random but different indices.
	var genIndex []int
	for len(genIndex) < numFlags {

		index := rand.Intn(len(flags))

		for i := 0; i < len(genIndex); i++ {

			if index == genIndex[i] {
				index = rand.Intn(len(flags))
				i = -1
			}
		}

		genIndex = append(genIndex, index)
	}

	// Add the corresponding flag of the previously generated
	// index to the string array "genFlags".
	var genFlags []string
	for i := 0; i < len(genIndex); i++ {
		genFlags = append(genFlags, flags[genIndex[i]])
	}

	// Generate final flag string.
	flagString := fmt.Sprintf("(%s)", strings.Join(genFlags, " "))

	return flagString, genFlags
}

// GenerateMsg returns a randomly generated message as
// second value and the message's byte length as first.
func GenerateMsg() (string, string) {

	// Choose mail version to generate.
	headerIndex := rand.Intn(5)

	// Generate number of lines of random strings to be
	// included in this message. 10 <= numLines <= 512.
	numLines := rand.Intn(503) + 10
	includeLines := make([]string, numLines)

	// Generate according number of lines.
	for i := 0; i < numLines; i++ {
		includeLines[i] = fmt.Sprintf("%s\r\n", GenerateString(64))
	}

	// Put together final message string.
	msg := fmt.Sprintf("%s%s", msgHeaders[headerIndex], strings.Join(includeLines, ""))
	msgLen := fmt.Sprintf("{%d}", len(msg))

	return msgLen, msg
}
