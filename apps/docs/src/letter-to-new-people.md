Hi there, Phoenix here.

One thing that seems to be not taught in school and mostly picked up on the fly is good software engineering practices. I don't even think SWE teaches these.

As part of the SH tech team, you are entrusted with the valuable task of making sure everything goes well. The frameworks you choose and the way you solve problems ultimately will not matter besides getting the expected functionality that everyone wants and making your deadlines. This means that, despite there
usually being more than one way to solve most of your problems, you're probably going to want to choose ones that will save you and your team members the most headache. There will also be people who may read your code long after you leave and go on to do greater things.

On top of this, some weeks you're a busy student without much time to get things done, and, well, may have to implement something the quick and dirty way and refactor later.
This is called "technical debt" and as a member of this team, please try to understand it very well.

I certainly think that being a part of this team will give you a lot of practical software engineering knowledge and the ability to have a positive impact on many people at a scale which you will have probably not had before: on the organizers, your team members, and ultimately, hundreds if not hopefully well over a thousand users. This document is to specifically give you advice and things to research about that will make you be able to work towards this end as best as possible and achieve amazing results. Following it may prevent a lot of unnecessary headaches, bad blood, and late nights.

### Everyone *will* make mistakes
In other words, nothing that you're going to write will be bug-less. Unexpected behavior will happen even after testing, and you have to ensure that things won't go
wrong when they do. For example, with sending emails in the email system, we made sure that both successes and failures were logged, and that they included the attempted contact email. What do you think would happen if we didn't do this? After a batch of decisions emails were sent, there potentially could have been hundreds of failures that would go unnoticed and lead to worse hackathon turnout. And if we had to re-do an email send, duplicate emails could have been sent if we didn't know who already probably received one. If we had to re-run the decisions algorithm too, then there would then be the potential for someone to receive both an acceptance and rejection email. Some simple future-proofing can avoid bad situations like these, which may I remind you can still happen even on well tested code.

### Atomicity
Noting the above example with decisions emails, what if we made the whole thing get kicked off with a single request to one API endpoint, deciding acceptances and *then* sending emails immediately after?

Well, what happens if all of the results created by the algorithm accidentally only accepted seniors? Congratulations, you would have sent a bunch of incorrect decisions emails and the organizers are now drafting an apology announcement (or maybe they put that on you). That's a big no no!

Preventing this situation would be done by breaking apart the functionality for the decisions feature into separately runnable pieces of code. Or in other words, making the decisions functionality more "atomic".

We did this by letting us calculate the decisions *first* using one endpoint, and after inspecting the data and confirming with the organizers that the statistics looked the way we wanted, we would hit another endpoint to queue up all of the emails which would then be sent out. Remember this when making features especially for public facing functionality.

### Resiliency
This is a more dev-ops related problem, but a fun one.

When working on SH you're going to be interacting with a lot of very valuable, mission critical data. But maybe someone was accidentally were testing something on the production database and deleted all of the precious user data. Maybe AWS is down (this actually happened on decision day lol). What now?

Well, you're probably thinking about backups. Great! But if your backups are online somewhere and the internet pipes stop working, then, hmm... what to do...
And are you only backing up production data? If so, how often should you do it? You're going to need a strategy.

My suggestion is to use a 3-2-1 backup strategy. If you find a better strategy, and everyone agrees with it, then great. use that. But at least strategize *how* you're going to do the backups. And the earlier you set that up the better. Basically:

* You want *3* copies of your data

* On *2* types of storage media (SSDs and HDDs can fail differently, use a tape drive if you want peak reliability)

* and most importantly *1* backup in a different location from your main db (offsite).

For our database we used Neon, which is just a managed Postgres db you can interface with. At the time of writing Neon allows you to rollback your database, but that couldn't be done past a certain elapsed time and was expensive. Thus, I wrote a script to dump the database to a personal server that I owned every 6 hours during the event, as well as before huge changes were made, like running the decisions algo. You will probably want something more substantial than that, and I hope this team manages to create a more permanent backup solution. Realistically, while people are signing up, a backup should be made at least once a day.

### Erroring, logging, and observability
Expect errors. Expect them in places where you think they would never happen, even though the program works with validated inputs or whatever. This is also one huge benefit for using Golang for our backend. Many functions return an error value and un-used variables force the program to not compile. That's an intentional choice by the authors of this language skewing developers to better programming practices.

Catching errors everywhere in your code can really save a lot of wasted time chasing odd behavior or bugs. You may also find that developing one feature months after another will cause some errors to throw that you've previously written. You'll be very glad you wrote them if that happens!

We also used something called a structured logging tool. This basically means that logs output to the console have a certain type and fields. 

This may seem like it just makes prettier console output, but if these are taken in by an observability tool (like Grafana) you can actually send out a notification to the tech team when a spike of errors hits the system, or when certain errors happen that you *know* shouldn't. Say for example, you could set an alert for a warning log saying that someone's UserID can't be found while running functionality on a page only authorized people have access to, potentially shedding light on a security vulnerability. That's big.

### Clever solution != a smart one
In other words, the very complicated cool way of doing things should be avoided over a simple and easy to understand solution.

You should aim to write software that looks so plain and simple that everything it does is obvious. Reduce complexity at all costs.

You might understand some complicated functionality very well. You may have put blood, sweat, and tears into it. But after you leave this team,
someone else may have to try to understand that, and if they can't, it might be easier for them to rewrite it. To prevent mass time loss across
you and other people, you should strive to make all functionality as clear and simple as possible. There are many guides online explaining how to do this.

Some hints I follow while programming that tells me when my code might be complicated are usually:

* my code begins nesting too much

* the order which functionality happens is not clear

* the functionality of a certain block of code is not clear

* what package is that function in again...?

I would also avoid

* tracking too much state (avoid state desyncronization)

And try to

* group functionality together as well as possible, file and folder wise
    * yes i have problems with the way the current backend is structured. look at a more mature go repository for the reason why.

### Your code should be self documenting
Comments are important, but they should not be over-utilized. Code should be simple enough to the point where what happens is clearly understandable, like I've been saying.

However, sometimes explaining what something does is best done with a comment. But when you do, understand that you are making an assumption that the person reading your code won't understand things too clearly. How much do you think they know?

I.e., you need to understand your audience when writing them.

In my opinion, people maintaining this codebase should have at least a working proficiency with the libraries and tech stack used, or are otherwise expected to learn enough to be at that level. Especially
whoever is reading the segment of code you're commenting on. When writing your comments, you should assume your readers know basics so you don't need to over-explain.

One practical example of this could be with HTTP status codes. People working on the HTTP API should be expected to know what these are, so you don't need to explain what 400 and 500 codes are and what kinds of errors they should be returned in.

However, this is to a point. When dealing with HTTP status codes, you could also assume that even those familiar with them may not have seen
certain ones in a while. Let's say you're returning a request on an error, and the user is unauthorized. We might just set a status code
field to the number 401, but to some people skimming through that http handler there is a good chance they may only have a vauge idea of what a 401 status is, and have to look it up on MDN web docs if they decide that's worth their time (key point: they might not). Instead, a better decision would be to name a variable called `StatusUnauthorized`,
set it to the value `401`, and then put that in the field.

Thankfully, Golang's net/http library has this for every status code. You should use them, like we did :)

### Communicating with organizers (or managers, or anyone else you have to report to)
I put this at the bottom because I probably don't want them to see this lol.

Please realize that while yes, the organizers may be quite pestering in wanting many things to be done, they have a lot of pressure
and unfortunately cannot just take your place and achieve everything they want. They have to instead entrust *you* and the entire tech
team things that they are responsible for. And the degree of that is super high, like making sure that everyone who shows up on the
day of the hackathon were the people who were actually accepted. Or sending out hundreds of notification emails. They certainly cannot do that by themselves.

Assuming the last people selected good organizers, they will probably be really ambitious and ask a lot from the tech team. That can be overwhelming.

Also, some of the things they ask for may also be unrealistic or incoherent with schedules, the state of the current codebase, etc.
While that might be intimidating at first, please understand that as a member of this team you still have a lot of agency with what gets
done and are completely free to argue about cutting features, suggesting an alternative solution to something they expect, etc. 

It is important that when they really want something done, to keep them in the loop and respond timely. But you are also able to set their expectations when you say what things are currently like. Please understand that you can use that to make the everyone on the team's life easier or more difficult. Not just yours.

### Blame
If something goes terribly wrong, I urge you to not try and single out any one person as the reason why something happened the way it happened.

Instead, let it be a learning opportunity for everyone. This can mean making a priority for testing, or peer reviewing code, spending more money to properly back up data and the like. Large mistakes in my opinion are usually caused by a large chain of events, and more safe guards or more attention to something can usually prevent mistakes in the future. But good organizational decisions have to be made.

Singling someone out is also going to de-moralize one out of the five or so developers you have and probably make them stressed as hell or not want to work, which can also put more responsibility on everyone else. Do you really want that? Please think calmly before releasing your frustration.

## Ending notes

I know this is a long read. But you also have a lot of work ahead of you. I hope you find some of the points here at least interesting or something worth some thought.
Again, it may save you and a lot of other people hassle they're probably never thinking about.

Good Luck Have Fun :)
