(function(qortexapi, $, undefined ) {
	qortexapi.rpc = function(methodName, input, callback) {
		var methodUrl = qortexapi.baseurl + "/" + methodName + ".json";
		var message = JSON.stringify(input);
		var req = $.ajax({
			type: "POST",
			url: methodUrl,
			contentType:"application/json; charset=utf-8",
			dataType:"json",
			processData: false,
			data: message
		});
		req.done(function(msg) {
			callback(msg, null);
		});
		req.fail(function(jqXHR, textStatus) {
			callback(null, textStatus);
		});
	};
})( window.qortexapi = window.qortexapi || {}, jQuery);



(function( qortexapi, undefined ) {
	qortexapi.Global = function() {
	}

	qortexapi.Global.prototype.GetAuthUserService = function(session) {
		var r = new qortexapi.AuthUserService()
		r.session = session
		return 
	}

	qortexapi.AuthUserService = function(){
	}

	qortexapi.AuthUserService.prototype.EditBroadcast = function(entryId, callback) {
		qortexapi.rpc("EditBroadcast", {"this": this, "params": {"entryId": entryId}}, callback)
	}


}( window.qortexapi = window.qortexapi || {} ));
