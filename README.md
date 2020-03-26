# Project K

Project K is a web-app for managing properties/tenants/maintenance requests like http://activebuilding.com. Our original goal is to provide an easy to use app for kos-kosan in Indonesia to manage their tenants, keep track of maintenance requests, and provide payment portal for easy rent transaction.

It took 3 months to build the entire app. It went from being a microservice architecture to a monolith due to costs, see https://medium.com/@obed.tandadjaja/reminder-to-self-always-start-with-a-monolith-e6dc1947982b for more information.

This specific backend is using http://gobuffalo.io as a framework. Some notes about Gobuffalo, I think the framework is very well done but it is not ready for production use yet. The framework seeks to create a Rails-like framework on Go, which is admirable, but I feel like it is trying to reinvent the wheel. People should just use Rails if they want built-in features and a community-tested web app framework, they should not make a similar thing in another language just for the sake of the language. Yes, Go is faster than Ruby (in some cases), but most of the time the reason that your web app is slow is not because of the language, it is because of your inefficient code. Overall, I think I was satisfied with Gobuffalo since it does give me out-of-the-box features and so I do not have to come up with my own solutions to common web app problems.

If I were to redo this project all over again, I would probably still use Go as the backend service. Yes, it is probably slower to develop on since it is much more verbose than Python and Ruby, but it is pleasant to code with. I will use a lighter web framework like Gin or Iris and use GORM on top of it for database accesses.

Anyway, the project ended in March 2020 since we did not find traction in demand of the product. Feel free to use this code or use it as an example project to learn more about Gobuffalo :)

## Database Setup

It looks like you chose to set up your application using a database! Fantastic!

The first thing you need to do is open up the "database.yml" file and edit it to use the correct usernames, passwords, hosts, etc... that are appropriate for your environment.

You will also need to make sure that **you** start/install the database of your choice. Buffalo **won't** install and start it for you.

### Create Your Databases

Ok, so you've edited the "database.yml" file and started your database, now Buffalo can create the databases in that file for you:

	$ buffalo pop create -a

## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

	$ buffalo dev

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.

**Congratulations!** You now have your Buffalo application up and running.

## What Next?

We recommend you heading over to [http://gobuffalo.io](http://gobuffalo.io) and reviewing all of the great documentation there.

Good luck!

[Powered by Buffalo](http://gobuffalo.io)
