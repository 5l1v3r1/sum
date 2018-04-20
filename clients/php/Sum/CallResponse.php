<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: proto/sum.proto

namespace Sum;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>sum.CallResponse</code>
 */
class CallResponse extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>bool success = 1;</code>
     */
    private $success = false;
    /**
     * Generated from protobuf field <code>string msg = 2;</code>
     */
    private $msg = '';
    /**
     * Generated from protobuf field <code>.sum.Data data = 3;</code>
     */
    private $data = null;

    public function __construct() {
        \GPBMetadata\Proto\Sum::initOnce();
        parent::__construct();
    }

    /**
     * Generated from protobuf field <code>bool success = 1;</code>
     * @return bool
     */
    public function getSuccess()
    {
        return $this->success;
    }

    /**
     * Generated from protobuf field <code>bool success = 1;</code>
     * @param bool $var
     * @return $this
     */
    public function setSuccess($var)
    {
        GPBUtil::checkBool($var);
        $this->success = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string msg = 2;</code>
     * @return string
     */
    public function getMsg()
    {
        return $this->msg;
    }

    /**
     * Generated from protobuf field <code>string msg = 2;</code>
     * @param string $var
     * @return $this
     */
    public function setMsg($var)
    {
        GPBUtil::checkString($var, True);
        $this->msg = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>.sum.Data data = 3;</code>
     * @return \Sum\Data
     */
    public function getData()
    {
        return $this->data;
    }

    /**
     * Generated from protobuf field <code>.sum.Data data = 3;</code>
     * @param \Sum\Data $var
     * @return $this
     */
    public function setData($var)
    {
        GPBUtil::checkMessage($var, \Sum\Data::class);
        $this->data = $var;

        return $this;
    }

}
