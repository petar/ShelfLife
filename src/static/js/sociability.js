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

}).call(this);
