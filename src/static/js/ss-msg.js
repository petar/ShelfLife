function ssAddMsg() {
	// social
	_.extend(ss.social, {

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
	});

	// view
	_.extend(ss.view, {
		MsgThread: Backbone.View.extend({

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
	});

	// ui
	_.extend(ss.view, {
		// ??
	});
}
