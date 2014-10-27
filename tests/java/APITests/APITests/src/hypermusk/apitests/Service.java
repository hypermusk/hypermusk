package hypermusk.apitests;


import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.Date;
import com.google.gson.annotations.SerializedName;
import com.google.gson.Gson;
import java.io.Serializable;

public class Service implements Serializable {





	
	// --- AuthorizeParams ---
	public static class AuthorizeParams {
		
	@SerializedName("Name")
	private String _name;



	public String getName() {
		return this._name;
	}
	public void setName(String _name) {
		this._name = _name;
	}



	}
	// --- AuthorizeResults ---
	public static class AuthorizeResults {
		
	@SerializedName("UseMap")
	private UseMap _useMap;

	@SerializedName("Err")
	private RemoteError _err;



	public UseMap getUseMap() {
		return this._useMap;
	}
	public void setUseMap(UseMap _useMap) {
		this._useMap = _useMap;
	}

	public RemoteError getErr() {
		return this._err;
	}
	public void setErr(RemoteError _err) {
		this._err = _err;
	}



	}
	// --- Authorize ---
	public AuthorizeResults authorize(String name) {
	
	AuthorizeResults results = null;
	AuthorizeParams params = new AuthorizeParams();
	params.setName(name);
	
	Api _api = Api.get();

	Gson gson = _api.gson();

	Map m = new HashMap();
	m.put("Params", params);
	m.put("This", this);

	String url = String.format("%s/Service/Authorize.json", _api.getBaseURL());
	String requestBody = gson.toJson(m);

	Api.RequestResult r = Api.request(url, requestBody, null);


	if (r.getErr() != null) {
		results = new AuthorizeResults();
		results.setErr(r.getErr());
		return results;
	}

	results = gson.fromJson(r.getReader(), AuthorizeResults.class);
	try {
		r.getReader().close();
	} catch (IOException e) {
	}

	

	return results;
	}


	
	// --- PermiessionDeniedParams ---
	public static class PermiessionDeniedParams {
		




	}
	// --- PermiessionDeniedResults ---
	public static class PermiessionDeniedResults {
		
	@SerializedName("Err")
	private RemoteError _err;



	public RemoteError getErr() {
		return this._err;
	}
	public void setErr(RemoteError _err) {
		this._err = _err;
	}



	}
	// --- PermiessionDenied ---
	public RemoteError permiessionDenied() {
	
	PermiessionDeniedResults results = null;
	PermiessionDeniedParams params = new PermiessionDeniedParams();
	
	Api _api = Api.get();

	Gson gson = _api.gson();

	Map m = new HashMap();
	m.put("Params", params);
	m.put("This", this);

	String url = String.format("%s/Service/PermiessionDenied.json", _api.getBaseURL());
	String requestBody = gson.toJson(m);

	Api.RequestResult r = Api.request(url, requestBody, null);


	if (r.getErr() != null) {
		results = new PermiessionDeniedResults();
		results.setErr(r.getErr());
		return results.getErr();
	}

	results = gson.fromJson(r.getReader(), PermiessionDeniedResults.class);
	try {
		r.getReader().close();
	} catch (IOException e) {
	}

	

	return results.getErr();
	}


	
	// --- GetReservedKeywordsForObjCParams ---
	public static class GetReservedKeywordsForObjCParams {
		




	}
	// --- GetReservedKeywordsForObjCResults ---
	public static class GetReservedKeywordsForObjCResults {
		
	@SerializedName("R")
	private ReservedKeywordsForObjC _r;

	@SerializedName("Err")
	private RemoteError _err;



	public ReservedKeywordsForObjC getR() {
		return this._r;
	}
	public void setR(ReservedKeywordsForObjC _r) {
		this._r = _r;
	}

	public RemoteError getErr() {
		return this._err;
	}
	public void setErr(RemoteError _err) {
		this._err = _err;
	}



	}
	// --- GetReservedKeywordsForObjC ---
	public GetReservedKeywordsForObjCResults getReservedKeywordsForObjC() {
	
	GetReservedKeywordsForObjCResults results = null;
	GetReservedKeywordsForObjCParams params = new GetReservedKeywordsForObjCParams();
	
	Api _api = Api.get();

	Gson gson = _api.gson();

	Map m = new HashMap();
	m.put("Params", params);
	m.put("This", this);

	String url = String.format("%s/Service/GetReservedKeywordsForObjC.json", _api.getBaseURL());
	String requestBody = gson.toJson(m);

	Api.RequestResult r = Api.request(url, requestBody, null);


	if (r.getErr() != null) {
		results = new GetReservedKeywordsForObjCResults();
		results.setErr(r.getErr());
		return results;
	}

	results = gson.fromJson(r.getReader(), GetReservedKeywordsForObjCResults.class);
	try {
		r.getReader().close();
	} catch (IOException e) {
	}

	

	return results;
	}


}

