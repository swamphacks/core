SwampHacks has a few major events which the tech team is required to provide services for. 

---
### Pre-applications release
* coming soon page
  * A landing page will be needed to inform people that SwampHacks will indeed be happening soon. More importantly however, this page allows us to begin constructing an email list of people
    who are interested early on, which we can send emails to later asking them to apply.
    A simple form with an email input will be needed. Make sure to validate emails appropriately (preferrably on both the front and back end) and return visible errors to the user.
  * Note: you should also try to use last year's emails in your email list.

### Applications Release
* This is when we open the forms for people to start registering for the hackathon. While MLH gives some basic guidelines for what information should be included (including the stuff that
they specifically need to know), the form data you include beyond this can help make deciding who to accept easier. 
* Note that while you are probably wanting as many people to apply as possible, you only will be letting a certain amount of people participate after decisions. More people signing up means more people will be
rejected, and this will ultimately effect a lot of people's self-perception of skill, career potential, etc. During XI we tried to keep things as fair as possible by also prioritizing people in "early-career," as well as
their passion for engineering on top of experience. We will get into this more later, but at the very least, you need to carefully consider what information will let you make the most informed, fair decisions later on.

Also, it might be a good idea to get data that allows you to track some statistics for getting an understanding of what kind of people are applying, as well as some stuff to compare to for future years.
Very vaugely, this could be data representative of people participating in the hackathon from in and out of state.

* Data that can show a person's experience
  * This is pretty simple. Have them upload their resume. We then stored the resume in an amazon S3 bucket (using cloudflare r2 as a wrapper), using their userID as the filename. Files should be able
    to replaced, preferrably under the last filename as before. Please make sure that only platform admins (the people who will be making decisions based off of resumes) and the uploaders
    are the only people able to see these resumes.
* Data that can show a person's passion
  * This is a bit trickier. We had people answer a couple of essay questions that we thought would get a good understanding on if people actually like what they do or are just in their major for the money. There's
    also some benefit in filtering out people who aren't willing to write something to apply.
* Email compliance laws
  * You should also include a check box & togglable setting for knowing whether or not certain candidates want emails from SwampHacks. This is to comply with laws that prevent buisnesses from spamming people.
    https://www.ftc.gov/business-guidance/resources/can-spam-act-compliance-guide-business

### Decisions

Decisions happened in XI like this:
* Have staff individually review each application, and assign each person an experience and passion score on a 1-5 scale
  * NOTE NOTE NOTE NOTE NOTE:
    * PLEAAASE do not just make a 1-5 button for each one. Instead (say for Experience) make 1 = "Never programmed before", 3 = "Has had an internship", 4 = "Has had multiple internships", 5 = "Has worked at FANNG or big tech".
      This way, *everyone* grading applications has a clear idea of what a 1, 2, 3, 4, and 5 mean, and aren't going to grade people differently. Represent the values 1-5 internally, but on the UI just show what they mean pls
* After all applications are done, run an algorithm that turns the raw scores into a ranked list, where every person is above and below someone else. Take the first X amount of values you want to hit your admissions quota.
  * We also accounted for certain quotas we had to hit, like having 40% acceptances be people in early career, and 60% of people be from UF. Alex Wala wrote a pretty good blog post about how we did it: 
* Send out acceptance and rejection letters to everyone
* If there are some rejected people you want part of the hackathon, change their user's status, and send them an email saying that they're actually accepted now
  * We didn't have a system in place for this, but because there are some people you are certain you may want accepted or rejected no matter what (perhaps exlcuding them from the quota used by the admissions algo too) you should make
    a feature for flaging people
* Now, for all of the accepted people, they will have to click a button to confirm their attendance. Not everyone will do this, so we need a way to let some of the rejected people back in. Which leads us to...

### Waitlist
When a person is rejected, they should be able to click a button somewhere that allows them to join a waitlist. The waitlist might grow a lot and become very big.

At some point before the event begins, you will need to let people off of the waitlist, accept them, and allow them to confirm their acceptance just like before. Again, not everyone who is switched to accepted will confirm.

The way you let people off of the waitlist is up to you. We did it in a rolling fashion in which, about a couple weeks before the hackathon, we would let 50 people off of the waitlist every couple of days, with the first people who signed
up for the waitlist being the first people to be let off (first-come-first-serve). Everyone who had not confirmed before a cycle was also automatically added to the back of the waitlist.

Again, when they are accepted, they get an email. As of XI we already have a system that handles this, but it is your job to test it and make sure that it works they way you and your team expects it to, and modify things as needed.

### Check-in & Redeemables
Now it's the day of the event. There are a few people lined up outside of a door, and many people just drove a long distance to get here. How are you going to make sure that anyone who hasn't actually signed up and been confirmed isn't let in?

On top of this, each person is going to get certain things like a t-shirt, and food throughout the event. We need to make sure we don't give out too many things and be able to track our inventory. Further, we might also want to know who attends what
workshops, espeically ones run by sponsor companies who might also want to know the specific people who attended. We call each t-shirt, food item, and workshop a "redeemable". When we check for them, we need to make sure the line also moves quickly.

For XI we solved this problem with a pretty cool badge system. Essentially, when you sign in you get a badge with a lanyard, that also has a thin NFC sticker on it. When people got redeemables we scan the nfc sticker (which is very very quick, quicker than aligning a camera for a QR-code)
and our database automatically updates, subtracting 1 from an item's counter and also updating a user to say that they recieved a redeemable.

NFC stickers are quite simple. They are a small circuit with a bit of memory (a stored unique identifier) that can be accessed by near-field radio. Most modern smartphones have support for this, but (at the time) API support on the web is very new, unlike RFID, a similar technology.

Our organizers thought that the RFID stickers looked too bulky, so they bought the NFC ones instead. Because there is very bad web support for NFC we had to build out a react native app for the scanning. That is why the swamphacks/xi-checkin-nfc repository exists.

Anyways, on to check-in

Since user's don't already have something like an NFC tag, QR-codes was our best option for quickly moving a line of people while having proof that people have been accepted and confirmed. We had the QR-code show up in the portal, but we also put them inside of emails sent out several days before the event. This
was done by storing the QR codes as images in S3, and putting their appropriate links inside the sent emails.

During the check-in line, people would show the QR-code to someone with a phone. The person would scan the QR-code, verify that the user's name and profile icon was theirs, be assigned a badge after scanning it, and on they went to hack.

However, there are some edge cases you should know:
 * Some people will arrive late
   * Have a simple late form people can fill out *before* arriving late, with their expected time and preferrably a Discord/phone number to contact. Make sure to have someone check them in once they arrive.
 * Some people will show up without having submitted an application, and we are under quota
   * Have them register on a form, accept them internally and check them in if everything looks good.
 * Some people will show up still on the waitlist / rejected, and we are under quota
   * If using the XI system, get the email associated with their Discord account, look them up in the database, change their status to accepted, and check them in. We should have a SQL script that does this btw
