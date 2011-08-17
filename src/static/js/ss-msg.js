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

	// model
	// view
	_.extend(ss.view, {
	});

	// ui
	_.extend(ss.view, {
		// ??
	});
}
