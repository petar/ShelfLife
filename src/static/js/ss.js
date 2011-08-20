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
				error: ecb,
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
				error: ecb,
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
				error: ecb,
				success: okcb,
				type: "GET",
				url: "/api/ss/Unlike"
			});
		},

		// followInfo
		followInfo: function(fid, okcb, ecb) {
			$.ajax({
				data: { "What": fid, },
				dataType: "json",
				error: ecb,
				success: function(data) { okcb(data.Count, data.Follows === "1"); },
				type: "GET",
				url: "/api/ss/FollowInfo"
			});
		},

		// setFollow
		setFollow: function(fid, okcb, ecb) {
			$.ajax({
				data: { "What": fid, },
				dataType: "json",
				error: ecb,
				success: okcb,
				type: "GET",
				url: "/api/ss/SetFollow"
			});
		},

		// unsetFollow
		unsetFollow: function(fid, okcb, ecb) {
			$.ajax({
				data: { "What": fid, },
				dataType: "json",
				error: ecb,
				success: okcb,
				type: "GET",
				url: "/api/ss/UnsetFollow"
			});
		},

		addMsg: function(attachTo, replyTo, body, okcb, ecb) {
			$.ajax({
				data: { "AttachTo": attachTo, "ReplyTo": replyTo, "Body": body },
				dataType: "json",
				error: ecb,
				success: okcb,
				type: "GET",
				url: "/api/ss/AddMsg"
			});
		},

		editMsg: function(msg, body, okcb, ecb) {
			$.ajax({
				data: { "Msg": msg, "Body": body },
				dataType: "json",
				error: ecb,
				success: okcb,
				type: "GET",
				url: "/api/ss/EditMsg"
			});
		},

		removeMsg: function(msg, okcb, ecb) {
			$.ajax({
				data: { "Msg": msg },
				dataType: "json",
				error: ecb,
				success: okcb,
				type: "GET",
				url: "/api/ss/RemoveMsg"
			});
		},

		findMsgAttachedTo: function(attachTo, okcb, ecb) {
			$.ajax({
				data: { "AttachTo": attachTo },
				dataType: "json",
				error: ecb,
				success: okcb,
				type: "GET",
				url: "/api/ss/FindMsgAttachedTo"
			});
		}
	};

	// ——— model-begin ——— 
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

		}),
	
		// MsgThread manages the data associated to a single message thread.
		// MsgThread triggers "change" when the whole set of messages has been fetched from
		// the server; it triggers "add" when a new message has been added, with one
		// argument carrying the ID of the msg
		MsgThread: Backbone.Model.extend({
		
			defaults: { "attachTo": null, "data": null, "fetched": false },

			// _msgIDMap maps message IDs to the index (in attribute 'data') of the
			// corresponding message
			_msgIDMap: {},

			initialize: function() { 
				_.bindAll(this, '_bringOK', '_addOK', '_removeOK');
			},

			// bring loads the messages in the thread from the server asynchronously;
			// when ready, it triggers 'change'
			bring: function() {
				var a = this.getAttachTo();
				if (!_.isString(a)) {
					console.log("fetching for msg thread without attach-to");
					return;
				}
				ss.social.findMsgAttachedTo(a, this._bringOK);
			},
			_bringOK: function(result) { 
				var d = result.Results;
				this._msgIDMap = {};
				_.each(d, function(m, i) { this._msgIDMap[m.id] = i }, this);
				this.set({"data": d, "fetched": true});
			},

			// add adds the given message; it triggers 'add' after asynchronously receiving
			// confirmation from the server that the addition was successful
			add: function(replyTo, body) {
				ss.social.addMsg(this.get("attachTo"), replyTo, body, this._addOK);
			},
			_addOK: function(r) {
				var msg = r.Msg;
				var d = this.get("data");
				var i = d.length;
				d.push(msg);
				this._msgIDMap[msg.id] = i;
				this.trigger("add", msg);
			},

			remove: function(msgID) {
				ss.social.removeMsg(msgID, _.bind(function() { this._removeOK(msgID); }, this));
			},
			_removeOK: function(msgID) {
				this.trigger("remove", msgID);
			},

			getAttachTo: function() { return this.get("attachTo"); },

			getLen: function() {
				var d = this.get("data");
				if (_.isArray(d)) {
					return d.length;
				} else {
					return 0;
				}
			},

			getByIndex: function(i) {
				var d = this.get("data");
				if (_.isArray(d)) {
					return d[i];
				} else {
					return undefined;
				}
			},

			getByID: function(x) { return this.get("data")[this._msgIDMap[x]]; },

			getAll: function() { return this.get("data"); }

		})
	};
	// ——— model-end ——— 

	// ——— view-begin ——— 
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

			events: { 'click': 'click' },

			initialize: function() {
				this.model = ss.vars.user;
				_.bindAll(this, 'render', 'refresh', 'update', 'click');
				this.model.bind('change', this.refresh);
				this.count = 0;
				this.like = false;
				this.refresh();
			},

			click: function() {
				var name = this.model.whoAmI();
				if (!_.isNull(name)) {
					if (this.like) {
						ss.social.unlike(this.options.fid, this.refresh, function() {});
					} else {
						ss.social.like(this.options.fid, this.refresh, function() {});
					}
				} else {
					ss.ui.showMustSignInBox();
				}
			},

			refresh: function() {
				ss.social.likeInfo(this.options.fid, this.update, function(){});
			},

			update: function(count, like) {
				this.count = count;
				this.like = like;
				this.render();
			},

			render: function() {
				$(this.el).html($("#ss-like-tmpl").tmpl());
				if (!this.like) {
					this.$("#wrap").addClass("like");
					this.$("#wrap").removeClass("unlike");
					this.$("#action").text("Like");
				} else {
					this.$("#wrap").addClass("unlike");
					this.$("#wrap").removeClass("like");
					this.$("#action").text("Unlike");
				}
				this.$("#footnote").text(this.count + " likes");
				return this;
			}
		}),

		// FollowButton is the view of the follow button
		// TODO: This has to extend LikeButton, not copy it
		FollowButton: Backbone.View.extend({

			tagName: "div",

			events: { 'click': 'click' },

			initialize: function() {
				this.model = ss.vars.user;
				_.bindAll(this, 'render', 'refresh', 'update', 'click');
				this.model.bind('change', this.refresh);
				this.count = 0;
				this.follow = false;
				this.refresh();
			},

			click: function() {
				var name = this.model.whoAmI();
				if (!_.isNull(name)) {
					if (this.follow) {
						ss.social.unsetFollow(this.options.fid, this.refresh, function() {});
					} else {
						ss.social.setFollow(this.options.fid, this.refresh, function() {});
					}
				} else {
					ss.ui.showMustSignInBox();
				}
			},

			refresh: function() {
				ss.social.followInfo(this.options.fid, this.update, function(){});
			},

			update: function(count, follow) {
				this.count = count;
				this.follow = follow;
				this.render();
			},

			render: function() {
				$(this.el).html($("#ss-follow-tmpl").tmpl());
				if (!this.follow) {
					this.$("#wrap").addClass("follow");
					this.$("#wrap").removeClass("unfollow");
					this.$("#action").text("Follow");
				} else {
					this.$("#wrap").addClass("unfollow");
					this.$("#wrap").removeClass("follow");
					this.$("#action").text("Unfollow");
				}
				this.$("#footnote").text(this.count + " followers");
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
		}),

		// MustSignInBox is an overlay UI view that informs the user they must sign in
		MustSignInBox: Backbone.View.extend({

			tagName: "div",

			events: { 'click #ok': 'ok' },

			initialize: function() {
				_.bindAll(this, 'ok', 'render');
				this.model = ss.vars.user;
			},
			
			render: function() {
				$(this.el).html($("#ss-mustsigninbox-tmpl").tmpl());
				return this;
			},

			ok: function() { 
				this.remove(); 
			},

			remove: function() {
				$(this.el).remove();
			}
		}),

		// MsgThread is the view that shows a message thread
		MsgThread: Backbone.View.extend({

			tagName: "div",
			
			events: { 
				'click span#post': '_onClickPost', 
				'click span#cancel': '_onClickCancel', 
				'click span#reply': '_onClickReply', 
				'click span#recancel': '_onClickReCancel', 
				'click #a-reply': '_onClickReplyLink',
				'click #a-remove': '_onClickRemoveLink',
				'click #b-remove': '_onClickRemoveLink',
			},

			model: null,

			// Must pass an 'attachTo' option
			initialize: function() {
				this.model = ss.misc.getMsgThreadModel(this.getAttachTo());
				_.bindAll(this, 'render', '_onChange', '_onAdd', '_onRemove', '_add');
				this.model.bind('change', this._onChange);
				this.model.bind('add', this._onAdd);
				this.model.bind('remove', this._onRemove);
			},

			getAttachTo: function() { return this.options.attachTo; },

			_onChange: function() { this.render(); },

			_onClickPost: function() { 
				var dText = this.$('.ss-msg-post textarea');
				this.model.add("", dText.val());
				this._initPostBox();
			},

			_onClickCancel: function() { this._initPostBox(); },

			_onClickReply: function(e) {
				var dBox = $(e.currentTarget).parent();
				var dText = $('textarea', dBox);
				this.model.add(dText.attr('replyTo'), dText.val());
				dBox.hide();
				this._initText(dText);
			},

			_onClickReCancel: function(e) { 
				var dBox = $(e.currentTarget).parent();
				var dText = $('textarea', dBox);
				dBox.hide();
				this._initText(dText);
			},

			_initPostBox: function() { this._initText(this.$('.ss-msg-post textarea')) },

			_initText: function(dText) { 
				dText.val('');
				dText.css('height', '30px');
				/* autoResize does not work when we use val() to change textarea
				 * contents.
				dText.autoResize({
					onResize : function() { $(this).css({opacity:0.8}); },
					animateCallback : function() { $(this).css({opacity:1}); },
					animateDuration : 200,
					extraSpace : 10,
					limit: 150
				});
				*/
			},

			_onAdd: function(m) { this._add(m); },

			_onRemove: function(msgID) {
				this.$('.ss-msg-root[msgID='+msgID+']').remove();
				this.$('.ss-msg-re[msgID='+msgID+']').remove();
			},

			/* msgJoin = { id, body, author_id, author_nym, attach, reply } */
			render: function() {
				$(this.el).html($("#ss-msg-thread-tmpl").tmpl());
				_.each(this.model.getAll(), this._add);
				this._initPostBox();
				return this;
			},

			_add: function(m) {
				var dThread = this.$('.ss-msg-thread');
				var dBox = $('.ss-msg-box', dThread);
				// Reply message
				if (m.reply != "") {
					var dRoot = $('.ss-msg-root[msgID='+m.reply+']', dBox);
					if (!_.isNull(dRoot) && !_.isUndefined(dRoot)) {
						var dReplies = $('div.ss-msg-replies', dRoot);
						dReplies.css("display", "block");
						var dReBox = $('.ss-msg-rebox', dReplies);
						var dRe = $("#ss-msg-re-tmpl").tmpl();
						dRe.attr('msgID', m.id);
						$('#b-remove', dRe).attr('msgID', m.id);
						$('.ss-msg-nym', dRe).text(m.author_nym);
						$('.ss-msg-info', dRe).text(m.modified);
						$('.ss-msg-body', dRe).text(m.body);
						dReBox.append(dRe);
					}
				} else { // Root level message
					var dRoot = $("#ss-msg-root-tmpl").tmpl();
					dRoot.attr('msgID', m.id);
					dRoot.attr('authorID', m.author_id);
					$('.ss-msg-nym', dRoot).text(m.author_nym);
					$('.ss-msg-info', dRoot).text(m.modified);
					$('.ss-msg-body', dRoot).text(m.body);
					$('.ss-msg-respond', dRoot).hide();
					$('#a-reply', dRoot).attr('msgID', m.id);
					$('#a-remove', dRoot).attr('msgID', m.id);
					dBox.append(dRoot);
					var dReplyText = $('textarea', dRoot);
					dReplyText.attr('replyTo', m.id);
					this._initText(dReplyText);
				}
			},

			_onClickReplyLink: function(e) {
				var msgID = $(e.currentTarget).attr('msgID');
				this.$('.ss-msg-root[msgID='+msgID+'] .ss-msg-respond').show();
			},

			_onClickRemoveLink: function(e) {
				var msgID = $(e.currentTarget).attr('msgID');
				this.model.remove(msgID);
			}

		}),
	};
	// ——— view-end ——— 

	ss.vars = {
		// user is the unique global User model, an instance of ss.model.User
		user: new ss.model.User,

		// _msgThreads is a private variable holding all message thread models
		_msgThreads: {},
	};

	ss.misc = {
		getMsgThreadModel: function(attachTo) {
			var r = ss.vars._msgThreads[attachTo];
			if (_.isUndefined(r)) {
				r = new ss.model.MsgThread({"attachTo": attachTo});
				ss.vars._msgThreads[attachTo] = r;
			}
			return r;
		}
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

		// showMustSignInBox shows UI that informs the user they must be signed in
		showMustSignInBox: function() {
			$("body").prepend((new ss.view.MustSignInBox()).render().el);
		},

		// showLikeButtons replaces all <div class="inject-like" fid="..."></div> with a
		// like button
		showLikeButtons: function() {
			_.each($(".inject-like"), function (q) {
				$(q).removeClass("inject-like");
				$(q).prepend((new ss.view.LikeButton({ "fid": $(q).attr("fid") })).render().el);
			});
		},

		// showFollowButtons replaces all <div class="inject-follow" fid="..."></div> with a
		// follow button
		showFollowButtons: function() {
			_.each($(".inject-follow"), function (q) {
				$(q).removeClass("inject-follow");
				$(q).prepend((new ss.view.FollowButton({ "fid": $(q).attr("fid") })).render().el);
			});
		},

		// showMessageThreads replaces all <div class="inject-msgthread" fid="..."></div> with a
		// message thread UI
		showMessageThreads: function() {
			_.each($(".inject-msgthread"), function (q) {
				var attachTo = $(q).attr("fid");
				$(q).removeClass("inject-msgthread");
				$(q).prepend((new ss.view.MsgThread({ "attachTo": attachTo })).render().el);
				ss.misc.getMsgThreadModel(attachTo).bring();
			});
		},
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
