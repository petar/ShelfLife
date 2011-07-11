(function() {

	// Save a reference to the global object
	var root = this;

	ss = {};

	// ss.util contains a mixed bag of utility functions like cookie manipulation, etc.
	ss.util = {

		// createCookie creates a new cookie that expires in given number of days
		// Code taken from http://www.quirksmode.org/js/cookies.html
		createCookie: function(name,value,days) {
			if (days) {
				var date = new Date();
				date.setTime(date.getTime()+(days*24*60*60*1000));
				var expires = "; expires="+date.toGMTString();
			} else var expires = "";
			document.cookie = name+"="+value+expires+"; path=/";
		},

		// readCookie reads a cookie with the given name.
		// Code taken from http://www.quirksmode.org/js/cookies.html
		readCookie: function(name) {
			var nameEQ = name + "=";
			var ca = document.cookie.split(';');
			for(var i=0;i < ca.length;i++) {
				var c = ca[i];
				while (c.charAt(0)==' ') c = c.substring(1,c.length);
				if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
			}
			return null;
		},

		// eraseCookie erases the cookie with the given name.
		// Code taken from http://www.quirksmode.org/js/cookies.html
		eraseCookie: function(name) {
			ss.util.createCookie(name, "", -1);
		}
	};

	// ss.login takes care of user sign up, sign in, sign out, username availability checks, etc.
	ss.login = {

		// signIn logs in the user with (login, password).
		signIn : function(login, password, okcb, errcb) {
			this.signOut();
			$.ajax({
				data: {
					"L": login,
					"P": password,
				},
				dataType: "json",
				error: errcb,
				success: okcb,
				type: "GET",
				url: "/api/ss/SignInLogin"
			});
		},

		// signOut logs out the currently logged in user, by removing the authentication
		// cookies.
		signOut : function() {
			ss.util.eraseCookie("SS-UserAuth");
			ss.util.eraseCookie("SS-UserInfo");
		},

		// whatIsMyName returns the name of the currently logged in person, or undefined
		// otherwise.
		whatIsMyName : function() {
			if (_.isNull(ss.util.readCookie("SS-UserAuth"))) return null; 
			return ss.util.readCookie("SS-UserInfo");
		},

		// signUp registers the user specified by the tuple (name, email, login, password).
		signUp : function(name, email, login, password, okcb, errcb) {
			$.ajax({
				data: {
					"N": name,
					"E": email,
					"L": login,
					"P": password,
				},
				dataType: "json",
				error: errcb,
				success: okcb,
				type: "GET",
				url: "/api/ss/SignUp"
			});
		},

		// isLoginAvailable check whether login is not already taken by a user.
		// okcb(isAvailable) takes a single bool argument.
		isLoginAvailable : function(login, okcb, errcb) {
			$.ajax({
				data: {
					"L": login,
				},
				dataType: "json",
				error: errcb,
				success: function(data) {
					okcb(data.Available === "1");
				},
				type: "GET",
				url: "/api/ss/IsLoginAvailable"
			});
		}
	};

	// ss.model contains backbone.js models for the various functionalities
	ss.model = {

		// User manages user authentication: sign-in, sign-out and sign-up
		User: Backbone.Model.extend({
		
			defaults: { "name": null },		  

			initialize: function() { this._refresh(); },

			_refresh: function() {
				this.set({"name": ss.login.whatIsMyName()});
			},

			// signIn tries to sign the (login,password)
			signIn: function(login, password, okcb, ecb) {
				ss.login.signOut();
				ss.login.signIn(login, password,
					_.bind(function(okcb) { 
						this._refresh(); 
						if (_.isFunction(okcb)) okcb(); 
					}, this, okcb),
					ecb);
			},

			// signOut logs out the currently logged in user
			signOut: function() {
				ss.login.signOut();
				this._refresh();
			},

			// signUp registers a new user with the backend system
			signUp: function(name, email, login, password, okcb, ecb) {
				ss.login.signUp(name, email, login, password, okcb, ecb);
			},

			whoAmI: function() {
				return ss.login.whatIsMyName();
			}

		})
	
	};

	// ss.view contains all backbone.js views for various UI elements
	ss.view = {
		
		// UserBar is the view of the navigation and user management bar on top of the page
		UserBar: Backbone.View.extend({

			tagName: "div",

			initialize: function() {
				this.model = ss.vars.user;
				_.bindAll(this, 'render');
				this.model.bind('change', this.render);
			},
			
			render: function() {
				$(this.el).html($("#ss-bar-tmpl").tmpl());
				var name = this.model.whoAmI();
				console.log(name);
				if (!_.isNull(name)) {
					this.$("#sign-in").css("display", "none");
					this.$("#sign-up").css("display", "none");
					this.$("#sign-name").text(name);
					this.$("#sign-name").css("display", "inline-block");
					this.$("#sign-out").css("display", "inline-block");
				} else {
					this.$("#sign-name").css("display", "none");
					this.$("#sign-out").css("display", "none");
					this.$("#sign-in").css("display", "inline-block");
					this.$("#sign-up").css("display", "inline-block");
				}
				return this;
			}
		}),

		SignInBox: Backbone.View.extend({

			tagName: "div",

			initialize: function() {
				this.model = ss.vars.user;
			},
			
			render: function() {
				$("#ss-signinbox-tmpl").tmpl().prependTo(this.el);
				return this;
			}
		})
	};

	ss.vars = {
		// user is the unique global User model, an instance of ss.model.User
		user: new ss.model.User
	};

	ss.ui = {
		// showUserBar inserts the user bar UI element inside el
		showUserBar: function(el) {
			$(el).prepend((new ss.view.UserBar()).render().el);
		},

		// showSignInBox shows UI that prompts the user to sign in
		showSignInBox: function() {
			$("body").prepend((new ss.view.SignInBox()).render().el);
		},

		// showSignUpBox shows UI that prompts the user to register
		showSignUpBox: function() {
		}
	};

}).call(this);
