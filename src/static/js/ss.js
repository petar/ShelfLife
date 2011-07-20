function _init(okcb, errcb) {

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
		},

		// loadTemplate loads jquery-style templates from the given URL and
		// inserts them into the DOM
		loadTemplate: function(url, okcb, errcb) {
			$.ajax({
				dataType: "html",
				error: errcb,
				success: function(data) {
					$('<div>').html(data).appendTo($('body'));
					if (_.isFunction(okcb)) okcb();
				},
				type: "GET",
				url: url
			});
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

	ss.social = {
		// likeInfo asynchronously returns the number of likes for an object identified
		// by the string fid, and whether the currently logge user (if any) likes this
		// object
		likeInfo: function(fid, okcb, ecb) {
			$.ajax({
				data: { "FID": fid, },
				dataType: "json",
				error: errcb,
				success: function(data) { okcb(data.Count, data.Likes === "1"); },
				type: "GET",
				url: "/api/ss/LikeInfo"
			});
		},

		// like records that the currently logged user likes object fid
		like: function(fid, okcb, ecb) {
			$.ajax({
				data: { "FID": fid, },
				dataType: "json",
				error: errcb,
				success: okcb,
				type: "GET",
				url: "/api/ss/Like"
			});
		},

		// unlike records that the currently logged user does not like object fid
		unlike: function(fid, okcb, ecb) {
			$.ajax({
				data: { "FID": fid, },
				dataType: "json",
				error: errcb,
				success: okcb,
				type: "GET",
				url: "/api/ss/Unlike"
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

		// SignInBox is an overlay UI view that handles the UX for a sign in
		SignInBox: Backbone.View.extend({

			tagName: "div",

			events: { 'click #ok': 'ok', 'click #cancel': 'cancel' },

			initialize: function() {
				_.bindAll(this, 'ok', 'cancel', 'render', '_signInOk', '_signInErr');
				this.model = ss.vars.user;
			},
			
			render: function() {
				$(this.el).html($("#ss-signinbox-tmpl").tmpl());
				return this;
			},

			ok: function() {
				if (this.busy()) return;
				this.busy(true);
				var u = this.$("#u").val()
				var p = this.$("#p").val()
				this.model.signIn(u, p, this._signInOk, this._signInErr);
			},

			cancel: function() { 
				this.busy(false);
				this.cancelled = true;
				this.remove(); 
			},

			busy: function(on) {
				if (_.isUndefined(on)) {
					return this._busy == true;
				}
				if (on) {
					this._busy = true;
					this.$('#ok').attr('disabled', true);
				} else {
					this._busy = false;
					this.$('#ok').removeAttr('disabled');
				}
			},

			_signInOk: function() {
				this.remove();
			},

			_signInErr: function() {
				this.busy(false);
				if (!this.cancelled) {
					this.$('#err0').show();
				}
			},

			remove: function() {
				$(this.el).remove();
			}
		}),

		// LikeButton is the view of the like button
		LikeButton: Backbone.View.extend({

			tagName: "div",

			initialize: function() {
				this.model = ss.vars.user;
				_.bindAll(this, 'render', 'update');
				this.model.bind('change', this.render);
				this.count = 0;
				this.like = false;
				ss.social.likeInfo(this.options.fid, this.update);
			},

			update: function(count, like) {
				this.count = count;
				this.like = like;
				this.render();
			},

			render: function() {
				$(this.el).html($("#ss-like-tmpl").tmpl());
				if (!this.like) {
					this.$("#action").addClass("like");
					this.$("#action").removeClass("unlike");
					this.$("#action").text("Like");
				} else {
					this.$("#action").addClass("unlike");
					this.$("#action").removeClass("like");
					this.$("#action").text("Unlike");
				}
				this.$("#footnote").text(this.count + " likes");
				return this;
			}
		}),

		// SignUpBox is an overlay UI view that handles the UX for a sign up
		// TODO: Add an overlay message acknowledging success
		SignUpBox: Backbone.View.extend({

			tagName: "div",

			events: { 'click #ok': 'ok', 'click #cancel': 'cancel' },

			initialize: function() {
				_.bindAll(this, 'ok', 'cancel', 'render', '_signUpOk', '_signUpErr');
				this.model = ss.vars.user;
			},
			
			render: function() {
				$(this.el).html($("#ss-signupbox-tmpl").tmpl());
				return this;
			},

			ok: function() {
				if (this.busy()) return;
				this.busy(true);
				var n = this.$("#n").val()
				var e = this.$("#e").val()
				var u = this.$("#u").val()
				var p = this.$("#p").val()
				this.model.signOut();
				this.model.signUp(n, e, u, p, this._signUpOk, this._signUpErr);
			},

			cancel: function() { 
				this.busy(false);
				this.cancelled = true;
				this.remove(); 
			},

			busy: function(on) {
				if (_.isUndefined(on)) {
					return this._busy == true;
				}
				if (on) {
					this._busy = true;
					this.$('#ok').attr('disabled', true);
				} else {
					this._busy = false;
					this.$('#ok').removeAttr('disabled');
				}
			},

			_signUpOk: function() {
				this.remove();
			},

			_signUpErr: function() {
				this.busy(false);
				if (!this.cancelled) {
					this.$('#err0').show();
				}
			},

			remove: function() {
				$(this.el).remove();
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
			$("body").prepend((new ss.view.SignUpBox()).render().el);
		},

		// showLikeButtons replaces all <div class="inject-button" fid="..."></div> with a
		// like button
		showLikeButtons: function() {
			_.each($(".inject-button"), function (q) {
				$(q).removeClass("inject-button");
				$(q).prepend((new ss.view.LikeButton({ "fid": $(q).attr("fid") })).render().el);
			});
		}
	};

	// load the UI templates
	ss.util.loadTemplate('/s/ui.x-jquery-tmpl', okcb, errcb);
}

// initSociability initializes the sociability framework.
// Since initialization requires some ajax calls, the respective callbacks
// are called when the initialization has truly completed.
function initSociability(okcb, errcb) {
	$(_init.call(this, okcb, errcb));
}
