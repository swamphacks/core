SwampHacks has a few major events which the tech team is required to provide
services for.

---
## Pre-applications release

Coming Soon Page

  * A landing page will be needed to inform people that SwampHacks will indeed
    be happening soon. More importantly however, this page allows us to begin
    constructing an email list of people who are interested early on, which we can
    send emails to later asking them to apply.

    A simple form with an email input will be needed. Make sure to validate
    emails appropriately (preferrably on both the front and back end) and return
    visible errors to the user.

  * Note: you should also try to use last year's emails (participants, people
    from the interest list, everyone) in your email list.

## Applications Release

This is when we open the forms for people to start registering for the
hackathon. While MLH gives some basic guidelines for what information should
be included (including the stuff that they specifically need to know), the form
data you include beyond this can help make deciding who to accept easier.

Note that while you are probably wanting as many people to apply as possible,
you only will be letting a certain number of people participate after decisions.
More people signing up means more people will be rejected, and this will
ultimately affect a lot of people's self-perception of skill, career potential,
etc. During XI we tried to keep things as fair as possible by also prioritizing
people in "early-career," as well as their passion for engineering on top of
experience. We will get into this more later, but at the very least, you need
to carefully consider what information will let you make the most informed, fair
decisions later on.

Also, it might be a good idea to get data that allows you to track some
statistics for getting an understanding of what kind of people are applying,
as well as some stuff to compare to for future years. For example, you could
get info on who is applying inside and outside of UF, and who is attempting
to join their first hackathon. 

  * Data that can show a person's experience

      * This is pretty simple. Have them upload their resume. We then stored the
        resume in an amazon S3 bucket (using cloudflare r2 as a wrapper), using their
        userID as the filename. Files should be able to be replaced, preferably under
        the last filename as before. Please make sure that only platform admins (the
        people who will be making decisions based off of resumes) and the uploaders are
        the only people able to see these resumes. For XI this was implemented with
        [pre-signed objects returned from S3]().
        # TODO: add link

  * Data that can show a person's passion

      * This is a bit trickier. We had people answer a couple of essay questions
        that we thought would get a good understanding on if people actually like what
        they do or are just in it for the money. There's also some benefit in
        filtering out people who aren't willing to write something to apply.
        A GPT-ed response should (...hopefully) be obvious as well.

  * Email compliance laws

      * You should also include a check box & togglable setting for
        knowing whether or not certain candidates want emails from SwampHacks.
        This is to comply with laws that prevent businesses from spamming people.
        https://www.ftc.gov/business-guidance/resources/can-spam-act-compliance-guide-bu
        siness

## Decisions

Decisions happened in XI like this:

  1. Have staff individually review each application, and assign each person an experience and passion score on a 1-5 scale

      * NOTE NOTE NOTE NOTE NOTE:

          * PLEAAASE do not just make a 1-5 button for each one. Instead (say for
            Experience) make 1 = "Never programmed before", 3 = "Has had an internship",
            4 = "Has had multiple internships", 5 = "Has worked at FAANG or big tech".
            This way, *everyone* grading applications have a clear idea of what a 1, 2,
            3, 4, and 5 mean, and aren't going to grade people differently. Represent the
            values 1-5 internally, but on the UI just show what they mean pls

  1. After all applications are done, run an algorithm that turns the raw scores
     into a ranked list, where every person is above and below someo  else. Take the
     first X amount of values you want to hit your admissions quota.

      * We also accounted for certain quotas we had to hit, like having 40%
        acceptances be people in early career, and 60% of people be from UF. Alex Wala
        wrote a [pretty good blog post]() about how we did it.
        # TODO: add link

  1. Send out acceptance and rejection letters to everyone

  1. If there are some rejected people you want part of the hackathon, change
     their user's status, and send them an email saying that they're actually
     accepted now

      * We didn't have a system in place for this, but because there are some
        people you are certain you may want accepted or rejected no matter what (perhaps
        excluding them from the quota used by the admissions algo too) you should make a
        feature for flagging people

  1. Now, for all of the accepted people, they will have to click a button to
     confirm their attendance. Not everyone will do this, so we need a way to let
     some of the rejected people back in. Which leads us to...

## Waitlist
When a person is rejected, they should be able to click a button somewhere that
allows them to join a waitlist. The waitlist might grow a lot and become very
big.

At some point before the event begins, you will need to let people off of the
waitlist, accept them, and allow them to confirm their acceptance just like
before. Again, not everyone who is accepted will confirm.

The way you let people off of the waitlist is up to you. We did it in a
rolling fashion in which, about a couple weeks before the hackathon, we
would let 50 people off of the waitlist every couple of days, with the first
people who signed up for the waitlist being the first people to be let off
(first-come-first-serve). Everyone who had not been confirmed before a cycle was
also automatically added to the back of the waitlist.

Again, when they are accepted, they get an email. As of XI we already have a
system that handles this, but it is your job to test it and make sure that it
works they way you and your team expects it to, and modify things as needed.

## Check-in & Redeemables
### Background
Now it's the day of the event. There are a few people lined up outside of a
door, and many people just drove a long distance to get here. How are you going
to make sure that anyone who hasn't actually signed up and been confirmed isn't
let in?

On top of this, each person is going to get certain things like a t-shirt,
and food throughout the event. We need to make sure we don't give out too many
things and be able to track our inventory. Further, we might also want to know
who attends what workshops, especially ones run by sponsor companies who might
also want to know the specific people who attended. We call each t-shirt, food
item, and workshop a "redeemable". When we check for them, we need to make sure
the line also moves quickly.

For XI we solved this problem with a pretty cool badge system. Essentially,
when you sign in you get a badge with a lanyard that also has a thin NFC
sticker on it. When people got redeemables we scan the nfc sticker (which is
very very quick, quicker than aligning a camera for a QR-code) and our database
automatically updates, subtracting 1 from an item's counter and also updating a
user to say that they received a redeemable.

NFC stickers are quite simple. They are a small circuit with a bit of memory (a
stored unique identifier) that can be accessed by near-field radio. Most modern
smartphones have support for this, but (at the time) API support on the web was
very new, unlike RFID, a similar technology.

Our organizers thought that the RFID stickers looked too bulky, so they
bought the NFC ones instead. Because there is very bad web support for NFC
we had to build out a react native app for the scanning. That is why the
swamphacks/xi-checkin-nfc repository exists.

### Anyways, on to check-in

Since users don't already have something like an NFC tag, QR-codes were our
best option for quickly moving a line of people while having proof that people
have been accepted and confirmed. We had the QR-code show up in the portal, but
we also put them inside of emails sent out several days before the event. This
was done by storing the QR codes as images in S3, and putting their appropriate
links inside the sent emails (if you do this again, please test this
rigorously).

During the check-in line, people would show the QR-code to someone with a phone.
The person would scan the QR-code, verify that the user's name and profile icon
was theirs, be assigned a badge after scanning it, and on they went to hack.

However, there are some edge cases you should know:

  * Some people will arrive late

      * Have a simple late form people can fill out *before* arriving late, with
        their expected time and preferably a Discord/phone number to contact. Make sure
        to have someone check them in once they arrive.

  * Some people will show up without having submitted an application, and we are
    under quota

      * Have them register on a form, accept them internally and check them in
        if everything looks good.

  * Some people will show up still on the waitlist / rejected, and we are under quota

      * If using the XI system, get the email associated with their Discord
        account, look them up in the database, change their status to accepted, and
        check them in. We should have a SQL script that does this btw

## Judging

### Process

The awards are organized like this:

  * You will have multiple "tracks" that people submit their projects for

      * One general one for everyone, as well as specific ones. The organizers
        are going to pick these out, and some can be in line with the hackathon's
        theme.

  * There will also be smaller awards, like "least vibe-coded" and "best game design".

For the general track, all of the judges will review a list of projects that
they are assigned according to a set scale. Projects should have multiple judges
look at them, and their scores should be combined via some algorithm.

For other tracks, the same process holds. But rather than have all of the judges
reviewing for every track, only the ones with the most relevant experience will
judge for a track.

To this end, in XI we broke up judging into two rounds: one general and one
track round. Judges had to review for both rounds but with the requirements for
general and the specific assigned track in mind, respectively.

How did they input their scores for each project? We used a platform:

---
 
### Judging platform
It's up to you if you want to make your own judging platform. That can also
be quite an in-depth project though, and hackathon judging is a problem
already solved by a good few repos out there. So for XI we used one called
[hackutd/jury](https://github.com/hackutd/jury?tab=readme-ov-file).

While its ranking system allowed for multiple prize types, its authentication system
was pretty bad. Multiple accounts had to be made for each judge for every prize
that they judged for, and during the second round we also ran into an issue that
prevented us from creating accounts that had to be quickly solved. You might
want to try out a different platform.

We also set it up pretty quickly for XI and weren't able to test it before the
event (yikes!)

Please test your judging platform before everything happens, please.

## After

Congratulations, at this point, the hackathon is over! Yippee! I hope that the
experience is overall something that teaches you a lot and is enjoyable. There
are a lot of moving parts, features, and testing, but you will come out as
a much more experienced engineer as a part of this team, and hopefully have some
stories to tell. A lot of people will have a direct benefit from the services
that the tech team provided, whether they know it or not.

If you want some advice on making sure that all of the major points of this
article goes well, please take a look at "Letter To New Tech Teams" in the main
section of this documentation.

Good Luck Have Fun :)
